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
	} else {
		log.Printf("TransactionHandler: Comment: '%s'", comment)
	}

	tx := database.Transaction{
		AccountID:     uint(accountID),
		CategoryID:    uint(categoryID),
		SubCategoryID: uint(subCategoryID),
		Type:          database.TypeTransaction(transactionType),
		Amount:        amount,
		Comment:       comment,
		Date:          time.Now(),
	}

	log.Printf("TransactionHandler: Creating transaction - AccountID=%d, CategoryID=%d, SubCategoryID=%d, Type=%d, Amount=%.2f", accountID, categoryID, subCategoryID, transactionType, amount)

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
	} else {
		log.Printf("UpdateTransactionHandler: Comment: '%s'", comment)
	}

	var originalTx database.Transaction
	if err := db.First(&originalTx, transactionID).Error; err != nil {
		log.Printf("UpdateTransactionHandler: Transaction with ID %d not found: %v", transactionID, err)
		http.Error(w, "transaction not found", http.StatusNotFound)
		return
	}

	dateStr := r.FormValue("Date")
	timeStr := r.FormValue("Time")

	newDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Printf("UpdateTransactionHandler: Invalid date '%s': %v", dateStr, err)
		http.Error(w, "problem with date. use format YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	var parsedTime time.Time
	if timeStr != "" {
		parsedTime, err = time.Parse("15:04:05", timeStr)
		if err != nil {
			log.Printf("UpdateTransactionHandler: Invalid time format '%s', falling back to original. Error: %v", timeStr, err)
			parsedTime = originalTx.Date
		}
	} else {
		log.Printf("UpdateTransactionHandler: Time not provided in form, using original transaction time.")
		parsedTime = originalTx.Date
	}

	transactionDate := time.Date(
		newDate.Year(),
		newDate.Month(),
		newDate.Day(),
		parsedTime.Hour(),
		parsedTime.Minute(),
		parsedTime.Second(),
		parsedTime.Nanosecond(),
		originalTx.Date.Location(),
	)

	log.Printf("UpdateTransactionHandler: Combined new date and time: %v", transactionDate)

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
