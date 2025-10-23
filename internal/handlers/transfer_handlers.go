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

func TransferHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	log.Printf("TransferHandler: Starting transfer processing")

	if r.Method != http.MethodPost {
		log.Printf("TransferHandler: Invalid method %s, redirecting to /", r.Method)
		http.Redirect(w, r, "/transactions", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Printf("TransferHandler: Error parsing form: %v", err)
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	log.Printf("TransferHandler: Form data received - %+v", r.Form)

	accountID, err := strconv.Atoi(r.FormValue("AccountID"))
	if err != nil {
		log.Printf("TransferHandler: Invalid account ID '%s': %v", r.FormValue("AccountID"), err)
		http.Error(w, "problem with account_id. use normal values", http.StatusBadRequest)
		return
	}

	if accountID < 0 {
		log.Printf("TransferHandler: Negative account ID %d", accountID)
		http.Error(w, "problem with id. use positive values", http.StatusBadRequest)
		return
	}

	transferAccountID, err := strconv.Atoi(r.FormValue("TransferAccountID"))
	if err != nil {
		log.Printf("TransferHandler: Invalid transfer account ID '%s': %v", r.FormValue("TransferAccountID"), err)
		http.Error(w, "problem with transferaccount_id. use normal values", http.StatusBadRequest)
		return
	}

	if transferAccountID < 0 {
		log.Printf("TransferHandler: Negative transfer account ID %d", transferAccountID)
		http.Error(w, "problem with id. use positive values", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue("Amount"), 64)
	if err != nil {
		log.Printf("TransferHandler: Invalid amount '%s': %v", r.FormValue("Amount"), err)
		http.Error(w, "problem with amount. use normal values", http.StatusBadRequest)
		return
	}

	comment := r.FormValue("Description")
	if comment == "" {
		log.Printf("TransferHandler: Comment is empty, will be stored as NULL")
		comment = " "
	} else {
		log.Printf("TransferHandler: Comment: '%s'", comment)
	}

	tx := database.TransferTransaction{
		AccountID:         uint(accountID),
		TransferAccountID: uint(transferAccountID),
		Type:              database.Transfer,
		Amount:            amount,
		Comment:           comment,
		Date:              time.Now(),
	}

	log.Printf("TransferHandler: Creating transfer - From AccountID=%d, To AccountID=%d, Amount=%.2f", accountID, transferAccountID, amount)

	err = database.TransferToAnotherAcc(db, tx)
	if err != nil {
		log.Printf("TransferHandler: Failed to process transfer: %v", err)
		http.Error(w, fmt.Sprintf("failed to transfer: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("TransferHandler: Transfer completed successfully - From AccountID=%d, To AccountID=%d, Amount=%.2f", accountID, transferAccountID, amount)
	http.Redirect(w, r, "/transactions", http.StatusSeeOther)
}

func UpdateTransferHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	log.Printf("UpdateTransferHandler: Starting transfer update")

	if r.Method != http.MethodPost {
		log.Printf("UpdateTransferHandler: Invalid method %s", r.Method)
		http.Redirect(w, r, "/transactions", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("UpdateTransferHandler: Error parsing form: %v", err)
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	log.Printf("UpdateTransferHandler: Form data received - %+v", r.Form)

	transferIDStr := r.FormValue("ID")
	if transferIDStr == "" {
		log.Printf("UpdateTransferHandler: Missing transfer ID")
		http.Error(w, "transfer ID is required", http.StatusBadRequest)
		return
	}

	transferID, err := strconv.Atoi(transferIDStr)
	if err != nil {
		log.Printf("UpdateTransferHandler: Invalid transfer ID '%s': %v", transferIDStr, err)
		http.Error(w, "problem with transfer id", http.StatusBadRequest)
		return
	}

	fromAccountID, err := strconv.Atoi(r.FormValue("AccountID"))
	if err != nil {
		log.Printf("UpdateTransferHandler: Invalid from account ID: %v", err)
		http.Error(w, "problem with from account id", http.StatusBadRequest)
		return
	}

	toAccountID, err := strconv.Atoi(r.FormValue("TransferAccountID"))
	if err != nil {
		log.Printf("UpdateTransferHandler: Invalid to account ID: %v", err)
		http.Error(w, "problem with to account id", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue("Amount"), 64)
	if err != nil {
		log.Printf("UpdateTransferHandler: Invalid amount: %v", err)
		http.Error(w, "problem with amount", http.StatusBadRequest)
		return
	}

	dateStr := r.FormValue("Date")
	log.Printf("UpdateTransferHandler: Raw date string: '%s'", dateStr)

	var transferDate time.Time
	if dateStr != "" {
		var err error
		location, _ := time.LoadLocation("Local")

		transferDate, err = time.ParseInLocation("2006-01-02T15:04", dateStr, location)
		if err != nil {
			log.Printf("UpdateTransferHandler: Failed to parse as datetime-local, trying other formats: %v", err)

			transferDate, err = time.ParseInLocation("2006-01-02 15:04", dateStr, location)
			if err != nil {
				log.Printf("UpdateTransferHandler: Failed to parse as '2006-01-02 15:04': %v", err)
				http.Error(w, "problem with date format. use YYYY-MM-DDTHH:MM", http.StatusBadRequest)
				return
			}
		}

		log.Printf("UpdateTransferHandler: Parsed date: %v", transferDate)
	} else {
		transferDate = time.Now()
		log.Printf("UpdateTransferHandler: Using current date: %v", transferDate)
	}

	description := r.FormValue("Description")

	var originalTransfer database.TransferTransaction
	if err := db.First(&originalTransfer, transferID).Error; err != nil {
		log.Printf("UpdateTransferHandler: Transfer with ID %d not found: %v", transferID, err)
		http.Error(w, "transfer not found", http.StatusNotFound)
		return
	}

	transfer := database.TransferTransaction{
		ID:                uint(transferID),
		AccountID:         uint(fromAccountID),
		TransferAccountID: uint(toAccountID),
		Amount:            amount,
		Comment:           description,
		Date:              transferDate,
	}

	err = database.UpdateTransfer(db, transfer)
	if err != nil {
		log.Printf("UpdateTransferHandler: Failed to update transfer: %v", err)
		http.Error(w, fmt.Sprintf("failed to update transfer: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("UpdateTransferHandler: Transfer ID %d updated successfully", transferID)
	http.Redirect(w, r, "/transactions", http.StatusSeeOther)
}

func DeleteTransferHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	log.Printf("DeleteTransferHandler: Starting transfer deletion")

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/transactions", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("DeleteTransferHandler: Error parsing form: %v", err)
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	transferID, err := strconv.Atoi(r.FormValue("ID"))
	if err != nil {
		log.Printf("DeleteTransferHandler: Invalid transfer ID: %v", err)
		http.Error(w, "invalid transfer ID", http.StatusBadRequest)
		return
	}

	err = database.DeleteTransfer(db, transferID)
	if err != nil {
		log.Printf("DeleteTransferHandler: Error deleting transfer: %v", err)
		http.Error(w, "error deleting transfer", http.StatusInternalServerError)
		return
	}

	log.Printf("DeleteTransferHandler: Transfer %d deleted successfully", transferID)
	http.Redirect(w, r, "/transactions", http.StatusSeeOther)
}
