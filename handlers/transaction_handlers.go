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

func TransactionHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
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

	categoryID, err := strconv.Atoi(r.FormValue("CategoryID"))
	if err != nil {
		http.Error(w, "problem with category_id. use normal values", http.StatusBadRequest)
		return
	}

	if categoryID < 0 {
		http.Error(w, "problem with category_id. use positive values", http.StatusBadRequest)
		return
	}

	subCategoryID, err := strconv.Atoi(r.FormValue("SubCategoryID"))
	if err != nil {
		http.Error(w, "problem with subcategory_id. use normal values", http.StatusBadRequest)
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

	tx := database.Transaction{
		AccountID:     uint(accountID),
		CategoryID:    uint(categoryID),
		SubCategoryID: uint(subCategoryID),
		Type:          database.TypeTransaction(transType),
		Amount:        amount,
		Comment:       comment,
		Date:          time.Now(),
	}

	err = database.AddTransaction(db, tx)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to add transaction: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("Transaction added and account balance updated")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
