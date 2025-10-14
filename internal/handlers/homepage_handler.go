package handlers

import (
	"log"
	"net/http"
	"text/template"

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
	accountsForView, err := getAccountsForView(db)
	if err != nil {
		log.Printf("NewTransactionPageHandler: Error getting accounts: %v", err)
		http.Error(w, "could not get accounts", http.StatusInternalServerError)
		return
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
