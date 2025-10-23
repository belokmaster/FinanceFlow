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

	var subCategoryID uint
	subCategoryIDStr := r.FormValue("SubCategoryID")
	if subCategoryIDStr != "" {
		subID, err := strconv.Atoi(subCategoryIDStr)
		if err != nil {
			log.Printf("TransactionHandler: Invalid subcategory ID '%s': %v", subCategoryIDStr, err)
			http.Error(w, "problem with subcategory_id. use normal values", http.StatusBadRequest)
			return
		}

		if subID < 0 {
			log.Printf("TransactionHandler: Negative subcategory ID %d", subID)
			http.Error(w, "problem with subcategory_id. use positive values", http.StatusBadRequest)
			return
		}
		subCategoryID = uint(subID)
		log.Printf("TransactionHandler: SubCategoryID selected: %d", subCategoryID)
	} else {
		log.Printf("TransactionHandler: SubCategoryID not selected")
	}

	transactionType, err := strconv.Atoi(r.FormValue("Type"))
	if err != nil {
		log.Printf("TransactionHandler: Invalid transaction type '%s': %v", r.FormValue("Type"), err)
		http.Error(w, "problem with transactionType. use normal values", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue("Amount"), 64)
	if err != nil {
		log.Printf("TransactionHandler: Invalid amount '%s': %v", r.FormValue("Amount"), err)
		http.Error(w, "problem with amount. use normal values", http.StatusBadRequest)
		return
	}

	comment := r.FormValue("Comment")

	if comment == "" {
		log.Printf("TransactionHandler: Comment is empty, will be stored as NULL")
		comment = " "
	} else {
		log.Printf("TransactionHandler: Comment: '%s'", comment)
	}

	dateStr := r.FormValue("Date")
	var transactionDate time.Time
	if dateStr != "" {
		location, locErr := time.LoadLocation("Local")
		if locErr != nil {
			location = time.UTC
		}

		// parsing form HTML datetime-local (YYYY-MM-DDTHH:MM)
		transactionDate, err = time.ParseInLocation("2006-01-02T15:04", dateStr, location)
		if err != nil {
			log.Printf("TransactionHandler: Invalid date format '%s': %v", dateStr, err)
			http.Error(w, "invalid date format", http.StatusBadRequest)
			return
		}
		log.Printf("TransactionHandler: User selected date: %v", transactionDate)
	} else {
		transactionDate = time.Now()
		log.Printf("TransactionHandler: Using current date: %v (Local: %v)", transactionDate, transactionDate.Local())
	}

	tx := database.Transaction{
		AccountID:     uint(accountID),
		CategoryID:    uint(categoryID),
		SubCategoryID: uint(subCategoryID),
		Type:          database.TypeTransaction(transactionType),
		Amount:        amount,
		Comment:       comment,
		Date:          transactionDate,
	}

	log.Printf("TransactionHandler: Creating transaction - AccountID=%d, CategoryID=%d, SubCategoryID=%d, Type=%d, Amount=%.2f, Date=%v", accountID, categoryID, subCategoryID, transactionType, amount, transactionDate)

	err = database.AddTransaction(db, tx)
	if err != nil {
		log.Printf("TransactionHandler: Failed to add transaction: %v", err)
		http.Error(w, fmt.Sprintf("failed to add transaction: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("TransactionHandler: Transaction added successfully - AccountID=%d, Type=%d, Amount=%.2f", accountID, transactionType, amount)
	http.Redirect(w, r, "/transactions", http.StatusSeeOther)
}

func DeleteTransactionHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	log.Printf("DeleteTransactionHandler: Starting transaction deletion")

	if r.Method != http.MethodPost {
		log.Printf("DeleteTransactionHandler: Invalid method %s, redirecting to /", r.Method)
		http.Redirect(w, r, "/transactions", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Printf("DeleteTransactionHandler: Error parsing form: %v", err)
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	log.Printf("DeleteTransactionHandler: Form data received - %+v", r.Form)

	transactionID, err := strconv.Atoi(r.FormValue("ID"))
	if err != nil {
		log.Printf("DeleteTransactionHandler: Invalid transaction ID '%s': %v", r.FormValue("ID"), err)
		http.Error(w, "problem with id. use normal values", http.StatusBadRequest)
		return
	}

	if transactionID < 0 {
		log.Printf("DeleteTransactionHandler: Negative transaction ID %d", transactionID)
		http.Error(w, "problem with id. use positive values", http.StatusBadRequest)
		return
	}

	log.Printf("DeleteTransactionHandler: Attempting to delete transaction ID %d", transactionID)

	err = database.DeleteTransaction(db, transactionID)
	if err != nil {
		log.Printf("DeleteTransactionHandler: Failed to delete transaction ID %d: %v", transactionID, err)
		http.Error(w, fmt.Sprintf("failed to delete transaction: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("DeleteTransactionHandler: Transaction ID %d deleted successfully", transactionID)
	http.Redirect(w, r, "/transactions", http.StatusSeeOther)
}

func UpdateTransactionHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	log.Printf("UpdateTransactionHandler: Starting transaction update")

	if r.Method != http.MethodPost {
		log.Printf("UpdateTransactionHandler: Invalid method %s, redirecting to /", r.Method)
		http.Redirect(w, r, "/transactions", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("UpdateTransactionHandler: Error parsing form: %v", err)
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	log.Printf("UpdateTransactionHandler: Form data received - %+v", r.Form)

	transactionID, err := strconv.Atoi(r.FormValue("ID"))
	if err != nil {
		log.Printf("UpdateTransactionHandler: Invalid transaction ID '%s': %v", r.FormValue("ID"), err)
		http.Error(w, "problem with transaction id. use normal values", http.StatusBadRequest)
		return
	}

	if transactionID < 0 {
		log.Printf("UpdateTransactionHandler: Negative transaction ID %d", transactionID)
		http.Error(w, "problem with transaction id. use positive values", http.StatusBadRequest)
		return
	}

	accountID, err := strconv.Atoi(r.FormValue("AccountID"))
	if err != nil {
		log.Printf("UpdateTransactionHandler: Invalid account ID '%s': %v", r.FormValue("AccountID"), err)
		http.Error(w, "problem with account_id. use normal values", http.StatusBadRequest)
		return
	}

	if accountID < 0 {
		log.Printf("UpdateTransactionHandler: Negative account ID %d", accountID)
		http.Error(w, "problem with account_id. use positive values", http.StatusBadRequest)
		return
	}

	categoryID, err := strconv.Atoi(r.FormValue("CategoryID"))
	if err != nil {
		log.Printf("UpdateTransactionHandler: Invalid category ID '%s': %v", r.FormValue("CategoryID"), err)
		http.Error(w, "problem with category_id. use normal values", http.StatusBadRequest)
		return
	}

	if categoryID < 0 {
		log.Printf("UpdateTransactionHandler: Negative category ID %d", categoryID)
		http.Error(w, "problem with category_id. use positive values", http.StatusBadRequest)
		return
	}

	var subCategoryID uint
	subCategoryIDStr := r.FormValue("SubCategoryID")
	if subCategoryIDStr != "" {
		subID, err := strconv.Atoi(subCategoryIDStr)
		if err != nil {
			log.Printf("UpdateTransactionHandler: Invalid subcategory ID '%s': %v", subCategoryIDStr, err)
			http.Error(w, "problem with subcategory_id. use normal values", http.StatusBadRequest)
			return
		}

		if subID < 0 {
			log.Printf("UpdateTransactionHandler: Negative subcategory ID %d", subID)
			http.Error(w, "problem with subcategory_id. use positive values", http.StatusBadRequest)
			return
		}
		subCategoryID = uint(subID)
		log.Printf("UpdateTransactionHandler: SubCategoryID selected: %d", subCategoryID)
	} else {
		log.Printf("UpdateTransactionHandler: SubCategoryID not selected")
	}

	transactionType, err := strconv.Atoi(r.FormValue("Type"))
	if err != nil {
		log.Printf("UpdateTransactionHandler: Invalid transaction type '%s': %v", r.FormValue("Type"), err)
		http.Error(w, "problem with transactionType. use normal values", http.StatusBadRequest)
		return
	}

	amount := r.FormValue("Amount")
	var newAmount float64

	if amount != "" {
		newAmount, err = strconv.ParseFloat(amount, 64)
		if err != nil {
			log.Printf("UpdateTransactionHandler: Invalid amount '%s': %v", r.FormValue("Amount"), err)
			http.Error(w, "problem with amount. use normal values", http.StatusBadRequest)
			return
		}
	}

	comment := r.FormValue("Description")
	if comment == "" {
		log.Printf("UpdateTransactionHandler: Comment is empty, will be stored as NULL")
		comment = " "
	} else {
		log.Printf("UpdateTransactionHandler: Comment: '%s'", comment)
	}

	var originalTx database.Transaction
	if err := db.First(&originalTx, transactionID).Error; err != nil {
		log.Printf("UpdateTransactionHandler: Transaction with ID %d not found: %v", transactionID, err)
		http.Error(w, "transaction not found", http.StatusNotFound)
		return
	}

	dateTimeStr := r.FormValue("Date")
	var transactionDate time.Time

	if dateTimeStr != "" {
		location, err := time.LoadLocation("Local")
		if err != nil {
			location = time.UTC
		}

		transactionDate, err = time.ParseInLocation("2006-01-02T15:04", dateTimeStr, location)
		if err != nil {
			log.Printf("UpdateTransactionHandler: Invalid datetime format '%s': %v", dateTimeStr, err)
			http.Error(w, "problem with datetime. use format YYYY-MM-DDTHH:MM", http.StatusBadRequest)
			return
		}
		log.Printf("UpdateTransactionHandler: User selected datetime: %v", transactionDate)
	} else {
		transactionDate = originalTx.Date
		log.Printf("UpdateTransactionHandler: Using original transaction date: %v (Local: %v)", transactionDate, transactionDate.Local())
	}

	tx := database.Transaction{
		ID:            uint(transactionID),
		AccountID:     uint(accountID),
		CategoryID:    uint(categoryID),
		SubCategoryID: subCategoryID,
		Type:          database.TypeTransaction(transactionType),
		Amount:        newAmount,
		Comment:       comment,
		Date:          transactionDate,
	}

	log.Printf("UpdateTransactionHandler: Updating transaction ID %d - AccountID=%d, CategoryID=%d, SubCategoryID=%d, Type=%d, Amount=%.2f, Date=%v",
		transactionID, accountID, categoryID, subCategoryID, transactionType, newAmount, transactionDate)

	err = database.UpdateTransaction(db, tx)
	if err != nil {
		log.Printf("UpdateTransactionHandler: Failed to update transaction ID %d: %v", transactionID, err)
		http.Error(w, fmt.Sprintf("failed to update transaction: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("UpdateTransactionHandler: Transaction ID %d updated successfully", transactionID)
	http.Redirect(w, r, "/transactions", http.StatusSeeOther)
}
