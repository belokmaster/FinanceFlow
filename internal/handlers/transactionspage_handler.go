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

	accounts, err := database.GetAccounts(db)
	if err != nil {
		log.Printf("NewTransactionPageHandler: Error getting accounts: %v", err)
		http.Error(w, "could not get accounts", http.StatusInternalServerError)
		return
	}

	var accountsForView []AccountView
	for _, acc := range accounts {
		symbol, _ := database.CurrencySymbols[acc.CurrencyCode]
		color := "#" + acc.Color
		iconFileName, _ := database.IconAccountFiles[acc.IconCode]
		iconHTML, _ := icons.AccountIconCache[iconFileName]

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
	log.Printf("NewTransactionPageHandler: Prepared %d accounts", len(accountsForView))

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
