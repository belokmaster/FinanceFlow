package handlers

import (
	"finance_flow/internal/database"
	"log"
	"net/http"
	"sort"
	"text/template"
	"time"

	"gorm.io/gorm"
)

func groupTransactionsByDate(transactions []TransactionView, transfers []TransferView) []GroupedTransactions {
	grouped := make(map[string]*GroupedTransactions)

	for _, t := range transactions {
		parsedDate, err := time.Parse("2006-01-02", t.Date)
		if err != nil {
			log.Printf("Error parsing date '%s': %v", t.Date, err)
			continue
		}
		date := parsedDate.Format("02.01.2006")

		if _, exists := grouped[date]; !exists {
			grouped[date] = &GroupedTransactions{
				Date:           date,
				TotalAmount:    0,
				CurrencySymbol: t.CurrencySymbol,
				Transactions:   []TransactionView{},
				Transfers:      []TransferView{},
			}
		}

		grouped[date].Transactions = append(grouped[date].Transactions, t)

		switch t.Type {
		case 0:
			grouped[date].TotalAmount += t.Amount
		case 1:
			grouped[date].TotalAmount -= t.Amount
		}
	}

	for _, t := range transfers {
		parsedDate, err := time.Parse("2006-01-02", t.Date)
		if err != nil {
			log.Printf("Error parsing date '%s': %v", t.Date, err)
			continue
		}
		date := parsedDate.Format("02.01.2006")

		if _, exists := grouped[date]; !exists {
			grouped[date] = &GroupedTransactions{
				Date:           date,
				TotalAmount:    0,
				CurrencySymbol: t.CurrencySymbol,
				Transactions:   []TransactionView{},
				Transfers:      []TransferView{},
			}
		}

		grouped[date].Transfers = append(grouped[date].Transfers, t)
	}

	var result []GroupedTransactions
	for _, group := range grouped {
		result = append(result, *group)
	}

	sort.Slice(result, func(i, j int) bool {
		dateI, _ := time.Parse("02.01.2006", result[i].Date)
		dateJ, _ := time.Parse("02.01.2006", result[j].Date)
		return dateI.After(dateJ)
	})

	return result
}

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

	log.Printf("NewTransactionPageHandler: Getting transfers from database")
	transfersForView, err := getTransfersForView(db)
	if err != nil {
		log.Printf("NewTransactionPageHandler: Error getting transfers: %v", err)
		http.Error(w, "could not get transfers", http.StatusInternalServerError)
		return
	}

	log.Printf("NewTransactionPageHandler: Retrieved %d transactions from database", len(transactionsForView))

	groupedTransactions := groupTransactionsByDate(transactionsForView, transfersForView)

	subcategoriesByParent := convertSubcategoriesToView(subCategories)
	categoriesForView := convertCategoriesToView(categories, subcategoriesByParent)

	pageData := TransactionPageData{
		Accounts:            accountsForView,
		Categories:          categoriesForView,
		Transactions:        transactionsForView,
		GroupedTransactions: groupedTransactions,
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
