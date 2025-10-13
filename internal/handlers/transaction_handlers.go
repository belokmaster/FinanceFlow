package handlers

import (
	"finance_flow/internal/database"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"
)

func TransactionHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	log.Printf("TransactionHandler: Starting transaction processing")

	if r.Method != http.MethodPost {
		log.Printf("TransactionHandler: Invalid method %s, redirecting to /", r.Method)
		http.Redirect(w, r, "/transactions", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("TransactionHandler: Error parsing form: %v", err)
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	log.Printf("TransactionHandler: Form data received - %+v", r.Form)

	accountID, err := strconv.Atoi(r.FormValue("AccountID"))
	if err != nil {
		log.Printf("TransactionHandler: Invalid account ID '%s': %v", r.FormValue("AccountID"), err)
		http.Error(w, "problem with account_id. use normal values", http.StatusBadRequest)
		return
	}

	if accountID < 0 {
		log.Printf("TransactionHandler: Negative account ID %d", accountID)
		http.Error(w, "problem with id. use positive values", http.StatusBadRequest)
		return
	}

	categoryID, err := strconv.Atoi(r.FormValue("CategoryID"))
	if err != nil {
		log.Printf("TransactionHandler: Invalid category ID '%s': %v", r.FormValue("CategoryID"), err)
		http.Error(w, "problem with category_id. use normal values", http.StatusBadRequest)
		return
	}

	if categoryID < 0 {
		log.Printf("TransactionHandler: Negative category ID %d", categoryID)
		http.Error(w, "problem with category_id. use positive values", http.StatusBadRequest)
		return
	}

	subCategoryID, err := strconv.Atoi(r.FormValue("SubCategoryID"))
	if err != nil {
		log.Printf("TransactionHandler: Invalid subcategory ID '%s': %v", r.FormValue("SubCategoryID"), err)
		http.Error(w, "problem with subcategory_id. use normal values", http.StatusBadRequest)
		return
	}

	transType, err := strconv.Atoi(r.FormValue("Type"))
	if err != nil {
		log.Printf("TransactionHandler: Invalid transaction type '%s': %v", r.FormValue("Type"), err)
		http.Error(w, "problem with transtype. use normal values", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue("Amount"), 64)
	if err != nil {
		log.Printf("TransactionHandler: Invalid amount '%s': %v", r.FormValue("Amount"), err)
		http.Error(w, "problem with amount. use normal values", http.StatusBadRequest)
		return
	}

	comment := r.FormValue("Comment")
	log.Printf("TransactionHandler: Comment: '%s'", comment)

	tx := database.Transaction{
		AccountID:     uint(accountID),
		CategoryID:    uint(categoryID),
		SubCategoryID: uint(subCategoryID),
		Type:          database.TypeTransaction(transType),
		Amount:        amount,
		Comment:       comment,
		Date:          time.Now(),
	}

	log.Printf("TransactionHandler: Creating transaction - AccountID=%d, CategoryID=%d, SubCategoryID=%d, Type=%d, Amount=%.2f", accountID, categoryID, subCategoryID, transType, amount)

	err = database.AddTransaction(db, tx)
	if err != nil {
		log.Printf("TransactionHandler: Failed to add transaction: %v", err)
		http.Error(w, fmt.Sprintf("failed to add transaction: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("TransactionHandler: Transaction added successfully - AccountID=%d, Type=%d, Amount=%.2f", accountID, transType, amount)
	http.Redirect(w, r, "/transactions", http.StatusSeeOther)
}
