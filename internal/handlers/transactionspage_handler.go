package handlers

import (
	"finance_flow/internal/database"
	"finance_flow/internal/icons"
	"log"
	"net/http"
	"text/template"

	"gorm.io/gorm"
)

// rechange it later. make more convienent
func NewTransactionPageHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB, path string) {
	log.Printf("NewTransactionPageHandler: Preparing data for new transaction page")

	if r.Method != http.MethodGet {
		log.Printf("NewTransactionPageHandler: Invalid method %s, redirecting to /", r.Method)
		http.Redirect(w, r, "/transactions", http.StatusSeeOther)
		return
	}

	log.Printf("NewTransactionPageHandler: Getting accounts from database")
	accounts, err := database.GetAccounts(db)
	if err != nil {
		log.Printf("NewTransactionPageHandler: Error getting accounts: %v", err)
		http.Error(w, "could not get accounts", http.StatusInternalServerError)
		return
	}

	log.Printf("NewTransactionPageHandler: Retrieved %d accounts from database", len(accounts))

	var accountsForView []AccountView
	for _, acc := range accounts {
		symbol, ok := database.CurrencySymbols[acc.CurrencyCode]
		if !ok {
			log.Printf("NewTransactionPageHandler: Unknown currency code %d for account %s, using default symbol", acc.CurrencyCode, acc.Name)
			symbol = "?"
		}

		color := "#" + acc.Color
		if color == "" {
			log.Printf("NewTransactionPageHandler: Empty color for account %s, using default", acc.Name)
			color = "#4cd67a"
		}

		iconFileName, ok := database.IconAccountFiles[acc.IconCode]
		if !ok {
			log.Printf("NewTransactionPageHandler: Unknown icon code %d for account %s, using default icon", acc.IconCode, acc.Name)
			iconFileName = "Wallet"
		}

		iconHTML, ok := icons.AccountIconCache[iconFileName]
		if !ok {
			log.Printf("NewTransactionPageHandler: Icon file %s not found in cache for account %s, using coin icon", iconFileName, acc.Name)
			iconHTML = icons.AccountIconCache["coin"]
		}

		accountsForView = append(accountsForView, AccountView{
			ID:             acc.ID,
			Name:           acc.Name,
			Balance:        acc.Balance,
			CurrencySymbol: symbol,
			Color:          color,
			IconKey:        iconFileName,
			IconHTML:       iconHTML,
		})
	}

	log.Printf("HomeHandler: Parsing template from path: %s", path)

	categories, err := database.GetCategories(db)
	if err != nil {
		log.Printf("NewTransactionPageHandler: Error getting categories: %v", err)
		http.Error(w, "could not get categories", http.StatusInternalServerError)
		return
	}

	subCategories, err := database.GetSubCategories(db)
	if err != nil {
		log.Printf("NewTransactionPageHandler: Error getting sub_categories: %v", err)
		http.Error(w, "could not get sub_categories", http.StatusInternalServerError)
		return
	}

	subcategoriesByParent := make(map[uint][]SubCategoryView)
	for _, subCat := range subCategories {
		color := "#" + subCat.Color
		iconFileName, _ := database.IconSubCategoryFiles[subCat.IconCode]
		iconHTML, _ := icons.SubCategoryIconCache[iconFileName]

		subcategoryView := SubCategoryView{
			ID:       subCat.ID,
			Name:     subCat.Name,
			Color:    color,
			IconKey:  iconFileName,
			IconHTML: iconHTML,
			ParentID: subCat.CategoryID,
		}
		subcategoriesByParent[subCat.CategoryID] = append(subcategoriesByParent[subCat.CategoryID], subcategoryView)
	}

	var categoriesForView []CategoryView
	for _, cat := range categories {
		color := "#" + cat.Color
		iconFileName, _ := database.IconCategoryFiles[cat.IconCode]
		iconHTML, _ := icons.CategoryIconCache[iconFileName]

		subcats := subcategoriesByParent[cat.ID]
		if subcats == nil {
			subcats = []SubCategoryView{}
		}

		categoriesForView = append(categoriesForView, CategoryView{
			ID:            cat.ID,
			Name:          cat.Name,
			Color:         color,
			IconKey:       iconFileName,
			IconHTML:      iconHTML,
			Subcategories: subcats,
		})
	}
	log.Printf("NewTransactionPageHandler: Prepared %d categories with their subcategories", len(categoriesForView))

	pageData := TransactionPageData{
		Accounts:   accountsForView,
		Categories: categoriesForView,
	}

	tmpl, err := template.ParseFiles(path)
	if err != nil {
		log.Printf("NewTransactionPageHandler: Error parsing template: %v", err)
		http.Error(w, "could not parse template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, pageData)
	if err != nil {
		log.Printf("NewTransactionPageHandler: Error executing template: %v", err)
	} else {
		log.Printf("NewTransactionPageHandler: Template executed successfully")
	}
}
