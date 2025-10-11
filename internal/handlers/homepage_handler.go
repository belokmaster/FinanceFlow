package handlers

import (
	"log"
	"net/http"
	"text/template"

	"finance_flow/internal/database"
	"finance_flow/internal/icons"

	"gorm.io/gorm"
)

func HomeHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB, path string) {
	log.Printf("HomeHandler: Processing home page request")

	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	log.Printf("HomeHandler: Getting accounts from database")
	accounts, err := database.GetAccounts(db)
	if err != nil {
		log.Printf("HomeHandler: Error getting accounts: %v", err)
		http.Error(w, "could not get accounts", http.StatusInternalServerError)
		return
	}

	log.Printf("HomeHandler: Retrieved %d accounts from database", len(accounts))

	var accountsForView []AccountView
	for _, acc := range accounts {
		symbol, ok := database.CurrencySymbols[acc.CurrencyCode]
		if !ok {
			log.Printf("HomeHandler: Unknown currency code %d for account %s, using default symbol", acc.CurrencyCode, acc.Name)
			symbol = "?"
		}

		color := "#" + acc.Color
		if color == "" {
			log.Printf("HomeHandler: Empty color for account %s, using default", acc.Name)
			color = "#4cd67a"
		}

		iconFileName, ok := database.IconAccountFiles[acc.IconCode]
		if !ok {
			log.Printf("HomeHandler: Unknown icon code %d for account %s, using default icon", acc.IconCode, acc.Name)
			iconFileName = "Wallet"
		}

		iconHTML, ok := icons.AccountIconCache[iconFileName]
		if !ok {
			log.Printf("HomeHandler: Icon file %s not found in cache for account %s, using coin icon", iconFileName, acc.Name)
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

	tmpl, err := template.ParseFiles(path)
	if err != nil {
		log.Printf("HomeHandler: Error parsing template: %v", err)
		http.Error(w, "could not parse template", http.StatusInternalServerError)
		return
	}

	pageData := HomePageData{
		Accounts: accountsForView,
		Icons:    icons.AccountIconCache,
	}

	err = tmpl.Execute(w, pageData)
	if err != nil {
		log.Printf("HomeHandler: Error executing template: %v", err)
		log.Printf("error while proccess working: %v", err)
	} else {
		log.Printf("HomeHandler: Template executed successfully, sent %d accounts to client", len(accountsForView))
	}
}
