package handlers

import (
	"finance_flow/internal/database"
	"log"
	"net/http"
	"text/template"

	"gorm.io/gorm"
)

func AnalyzePageHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB, path string) {
	log.Printf("AnalyzePageHandler: Preparing data for new transaction page")

	if r.Method != http.MethodGet {
		log.Printf("AnalyzePageHandler: Invalid method %s, redirecting to /", r.Method)
		http.Redirect(w, r, "/analyze", http.StatusSeeOther)
		return
	}

	log.Printf("AnalyzePageHandler: Getting accounts from database")
	accountsForView, err := getAccountsForView(db)
	if err != nil {
		log.Printf("AnalyzePageHandler: Error getting accounts: %v", err)
		http.Error(w, "could not get accounts", http.StatusInternalServerError)
		return
	}

	log.Printf("AnalyzePageHandler: Getting categories from database")
	categories, err := database.GetCategories(db)
	if err != nil {
		log.Printf("AnalyzePageHandler: Error getting categories: %v", err)
		http.Error(w, "could not get categories", http.StatusInternalServerError)
		return
	}

	log.Printf("AnalyzePageHandler: Retrieved %d categories from database", len(categories))

	log.Printf("NewTransacAnalyzePageHandlertionPageHandler: Getting subcategories from database")
	subCategories, err := database.GetSubCategories(db)
	if err != nil {
		log.Printf("AnalyzePageHandler: Error getting sub_categories: %v", err)
		http.Error(w, "could not get sub_categories", http.StatusInternalServerError)
		return
	}

	log.Printf("AnalyzePageHandler: Retrieved %d sub_categories from database", len(subCategories))

	log.Printf("AnalyzePageHandler: Getting transactions from database")
	transactionsForView, err := getTransactionsForView(db)
	if err != nil {
		log.Printf("AnalyzePageHandler: Error getting transactions: %v", err)
		http.Error(w, "could not get transactions", http.StatusInternalServerError)
		return
	}

	log.Printf("AnalyzePageHandler: Retrieved %d transactions from database", len(transactionsForView))

	log.Printf("AnalyzePageHandler: Getting transfers from database")
	transfersForView, err := getTransfersForView(db)
	if err != nil {
		log.Printf("AnalyzePageHandler: Error getting transfers: %v", err)
		http.Error(w, "could not get transfers", http.StatusInternalServerError)
		return
	}

	log.Printf("AnalyzePageHandler: Retrieved %d transactions from database", len(transactionsForView))

	groupedTransactions := groupTransactionsByDate(transactionsForView, transfersForView)

	subcategoriesByParent := convertSubcategoriesToView(subCategories)
	categoriesForView := convertCategoriesToView(categories, subcategoriesByParent)

	pageData := AnalyzePageData{
		Accounts:            accountsForView,
		Categories:          categoriesForView,
		Transactions:        transactionsForView,
		GroupedTransactions: groupedTransactions,
	}

	log.Printf("AnalyzePageHandler: Parsing template from path: %s", path)
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		log.Printf("AnalyzePageHandler: Error parsing template: %v", err)
		http.Error(w, "could not parse template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, pageData)
	if err != nil {
		log.Printf("AnalyzePageHandler: Error executing template: %v", err)
	} else {
		log.Printf("AnalyzePageHandler: Template executed successfully")
	}
}
