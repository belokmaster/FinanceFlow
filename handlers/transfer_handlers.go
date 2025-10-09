package handlers

import (
	"finance_flow/database"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"
)

func TransferHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	accountID, err := strconv.Atoi(r.FormValue("AccountID"))
	if err != nil {
		http.Error(w, "problem with account_id. use normal values", http.StatusBadRequest)
		return
	}

	if accountID < 0 {
		http.Error(w, "problem with id. use positive values", http.StatusBadRequest)
		return
	}

	transferAccountID, err := strconv.Atoi(r.FormValue("TransferAccountID"))
	if err != nil {
		http.Error(w, "problem with transferaccount_id. use normal values", http.StatusBadRequest)
		return
	}

	if transferAccountID < 0 {
		http.Error(w, "problem with id. use positive values", http.StatusBadRequest)
		return
	}

	transType, err := strconv.Atoi(r.FormValue("Type"))
	if err != nil {
		http.Error(w, "problem with transtype. use normal values", http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue("Amount"), 64)
	if err != nil {
		http.Error(w, "problem with amount. use normal values", http.StatusBadRequest)
		return
	}

	comment := r.FormValue("Comment")

	tx := database.TransferTransaction{
		AccountID:         uint(accountID),
		TransferAccountID: uint(transferAccountID),
		Type:              database.TypeTransaction(transType),
		Amount:            amount,
		Comment:           comment,
		Date:              time.Now(),
	}

	err = database.TransferToAnotherAcc(db, tx)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to transfer: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("Transfer added and accounts balance updated")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
