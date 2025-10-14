package handlers

import (
	"finance_flow/internal/database"
	"log"
	"net/http"
	"text/template"

	"gorm.io/gorm"
)

func NewTransactionPageHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB, path string) {
	log.Printf("NewTransactionPageHandler: Preparing data for new transaction page")

	if r.Method != http.MethodGet {
		log.Printf("NewTransactionPageHandler: Invalid method %s, redirecting to /", r.Method)
		http.Redirect(w, r, "/transactions", http.StatusSeeOther)
		return
	}

	log.Printf("NewTransactionPageHandler: Getting accounts from database")
	accountsForView, err := getAccountsForView(db)
	if err != nil {
		log.Printf("NewTransactionPageHandler: Error getting accounts: %v", err)
		http.Error(w, "could not get accounts", http.StatusInternalServerError)
		return
	}

	log.Printf("NewTransactionPageHandler: Getting categories from database")
	categories, err := database.GetCategories(db)
	if err != nil {
		log.Printf("NewTransactionPageHandler: Error getting categories: %v", err)
		http.Error(w, "could not get categories", http.StatusInternalServerError)
		return
	}

	log.Printf("NewTransactionPageHandler: Retrieved %d categories from database", len(categories))

	log.Printf("NewTransactionPageHandler: Getting subcategories from database")
	subCategories, err := database.GetSubCategories(db)
	if err != nil {
		log.Printf("NewTransactionPageHandler: Error getting sub_categories: %v", err)
		http.Error(w, "could not get sub_categories", http.StatusInternalServerError)
		return
	}

	log.Printf("NewTransactionPageHandler: Retrieved %d sub_categories from database", len(subCategories))

	log.Printf("NewTransactionPageHandler: Getting transactions from database")
	transactionsForView, err := getTransactionsForView(db)
	if err != nil {
		log.Printf("NewTransactionPageHandler: Error getting transactions: %v", err)
		http.Error(w, "could not get transactions", http.StatusInternalServerError)
		return
	}

	log.Printf("NewTransactionPageHandler: Retrieved %d transactions from database", len(transactionsForView))

	subcategoriesByParent := convertSubcategoriesToView(subCategories)
	categoriesForView := convertCategoriesToView(categories, subcategoriesByParent)

	pageData := TransactionPageData{
		Accounts:     accountsForView,
		Categories:   categoriesForView,
		Transactions: transactionsForView,
	}

	log.Printf("NewTransactionPageHandler: Parsing template from path: %s", path)
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
