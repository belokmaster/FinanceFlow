package handlers

import (
	"encoding/json"
	"finance_flow/internal/database"
	"fmt"
	htmlTemplate "html/template"
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

	type analyzeTransaction struct {
		Type         int     `json:"type"`
		Amount       float64 `json:"amount"`
		CategoryName string  `json:"categoryName"`
		Date         string  `json:"date"`
		Currency     string  `json:"currency"`
	}

	type analyzeCategory struct {
		Name     string `json:"name"`
		Color    string `json:"color"`
		IconHTML string `json:"iconHtml"`
	}

	transactionsPayload := make([]analyzeTransaction, 0, len(transactionsForView))
	for _, tx := range transactionsForView {
		transactionsPayload = append(transactionsPayload, analyzeTransaction{
			Type:         int(tx.Type),
			Amount:       tx.Amount,
			CategoryName: tx.CategoryName,
			Date:         tx.Date.Format("2006-01-02"),
			Currency:     tx.CurrencySymbol,
		})
	}

	categoryMetaMap := make(map[string]analyzeCategory)
	for _, cat := range categoriesForView {
		categoryMetaMap[cat.Name] = analyzeCategory{
			Name:     cat.Name,
			Color:    fmt.Sprintf("#%s", cat.Color),
			IconHTML: string(cat.IconHTML),
		}
	}

	categoriesPayload := make([]analyzeCategory, 0, len(categoryMetaMap))
	for _, meta := range categoryMetaMap {
		categoriesPayload = append(categoriesPayload, meta)
	}

	transactionsJSONBytes, err := json.Marshal(transactionsPayload)
	if err != nil {
		log.Printf("AnalyzePageHandler: Error marshaling transactions payload: %v", err)
		http.Error(w, "could not build analytics payload", http.StatusInternalServerError)
		return
	}

	categoriesJSONBytes, err := json.Marshal(categoriesPayload)
	if err != nil {
		log.Printf("AnalyzePageHandler: Error marshaling categories payload: %v", err)
		http.Error(w, "could not build analytics payload", http.StatusInternalServerError)
		return
	}

	pageData := AnalyzePageData{
		Accounts:            accountsForView,
		Categories:          categoriesForView,
		Transactions:        transactionsForView,
		GroupedTransactions: groupedTransactions,
		TransactionsJSON:    htmlTemplate.JS(transactionsJSONBytes),
		CategoriesJSON:      htmlTemplate.JS(categoriesJSONBytes),
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
