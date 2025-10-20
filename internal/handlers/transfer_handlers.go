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
		http.Redirect(w, r, "/transactions", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("UpdateTransferHandler: Error parsing form: %v", err)
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	transferID, err := strconv.Atoi(r.FormValue("ID"))
	if err != nil {
		log.Printf("UpdateTransferHandler: Invalid transfer ID: %v", err)
		http.Error(w, "invalid transfer ID", http.StatusBadRequest)
		return
	}

	fromAccountID, err := strconv.Atoi(r.FormValue("AccountID"))
	if err != nil {
		log.Printf("UpdateTransferHandler: Invalid from account ID: %v", err)
		http.Error(w, "invalid from account ID", http.StatusBadRequest)
		return
	}

	toAccountID, err := strconv.Atoi(r.FormValue("TransferAccountID"))
	if err != nil {
		log.Printf("UpdateTransferHandler: Invalid to account ID: %v", err)
		http.Error(w, "invalid to account ID", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue("Amount"), 64)
	if err != nil {
		log.Printf("UpdateTransferHandler: Invalid amount: %v", err)
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	description := r.FormValue("Description")
	dateStr := r.FormValue("Date")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Printf("UpdateTransferHandler: Invalid date: %v", err)
		http.Error(w, "invalid date", http.StatusBadRequest)
		return
	}

	transfer := database.TransferTransaction{
		ID:                uint(transferID),
		AccountID:         uint(fromAccountID),
		TransferAccountID: uint(toAccountID),
		Amount:            amount,
		Comment:           description,
		Date:              date,
		Type:              database.Transfer,
	}

	err = database.UpdateTransfer(db, transfer)
	if err != nil {
		log.Printf("UpdateTransferHandler: Error updating transfer: %v", err)
		http.Error(w, "error updating transfer", http.StatusInternalServerError)
		return
	}

	log.Printf("UpdateTransferHandler: Transfer %d updated successfully", transferID)
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
