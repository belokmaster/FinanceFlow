package handlers

import (
	"finance_flow/database"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

func CreateAccountHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	accountID, err := strconv.Atoi(r.FormValue("ID"))
	if err != nil {
		http.Error(w, "problem with id. use normal values", http.StatusBadRequest)
		return
	}

	if accountID < 0 {
		http.Error(w, "problem with id. use positive values", http.StatusBadRequest)
		return
	}

	name := r.FormValue("Name")
	name = strings.TrimSpace(name) // why it works????????

	currency_id, err := strconv.Atoi(r.FormValue("Currency"))
	if err != nil {
		http.Error(w, "problem with currency. use normal values", http.StatusBadRequest)
		return
	}

	if currency_id < 0 {
		http.Error(w, "problem with currency. use positive values", http.StatusBadRequest)
		return
	}

	balance, err := strconv.ParseFloat(r.FormValue("Balance"), 64)
	if err != nil {
		http.Error(w, "problem with balance. use normal values", http.StatusBadRequest)
		return
	}

	acc := database.Account{
		Name:         name,
		Balance:      balance,
		CurrencyCode: database.TypeCurrency(currency_id),
	}

	err = database.CreateNewAccount(db, acc)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create a new acc: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("A new account was created sucessfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteAccountHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	accountID, err := strconv.Atoi(r.FormValue("ID"))
	if err != nil {
		http.Error(w, "problem with id. use normal values", http.StatusBadRequest)
		return
	}

	if accountID < 0 {
		http.Error(w, "problem with id. use positive values", http.StatusBadRequest)
		return
	}

	err = database.DeleteAccount(db, accountID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to delete a acc: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("A account was deleted sucessfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
