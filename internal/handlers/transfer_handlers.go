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
		http.Redirect(w, r, "/", http.StatusSeeOther)
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

	transType, err := strconv.Atoi(r.FormValue("Type"))
	if err != nil {
		log.Printf("TransferHandler: Invalid transaction type '%s': %v", r.FormValue("Type"), err)
		http.Error(w, "problem with transtype. use normal values", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue("Amount"), 64)
	if err != nil {
		log.Printf("TransferHandler: Invalid amount '%s': %v", r.FormValue("Amount"), err)
		http.Error(w, "problem with amount. use normal values", http.StatusBadRequest)
		return
	}

	comment := r.FormValue("Comment")
	log.Printf("TransferHandler: Comment: '%s'", comment)

	tx := database.TransferTransaction{
		AccountID:         uint(accountID),
		TransferAccountID: uint(transferAccountID),
		Type:              database.TypeTransaction(transType),
		Amount:            amount,
		Comment:           comment,
		Date:              time.Now(),
	}

	log.Printf("TransferHandler: Creating transfer - From AccountID=%d, To AccountID=%d, Type=%d, Amount=%.2f", accountID, transferAccountID, transType, amount)

	err = database.TransferToAnotherAcc(db, tx)
	if err != nil {
		log.Printf("TransferHandler: Failed to process transfer: %v", err)
		http.Error(w, fmt.Sprintf("failed to transfer: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("TransferHandler: Transfer completed successfully - From AccountID=%d, To AccountID=%d, Amount=%.2f", accountID, transferAccountID, amount)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
