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

	balance, err := strconv.ParseFloat(r.FormValue("Balance"), 64)
	if err != nil {
		http.Error(w, "problem with balance. use normal values", http.StatusBadRequest)
		return
	}

	currency_id, err := strconv.Atoi(r.FormValue("Currency"))
	if err != nil {
		http.Error(w, "problem with currency. use normal values", http.StatusBadRequest)
		return
	}

	if currency_id < 0 {
		http.Error(w, "problem with currency. use positive values", http.StatusBadRequest)
		return
	}

	color := r.FormValue("Color")

	icon_id, err := strconv.Atoi(r.FormValue("Icon"))
	if err != nil {
		http.Error(w, "problem with icon. use normal values", http.StatusBadRequest)
		return
	}

	if currency_id < 0 {
		http.Error(w, "problem with icon. use positive values", http.StatusBadRequest)
		return
	}

	acc := database.Account{
		Name:         name,
		Balance:      balance,
		CurrencyCode: database.TypeCurrency(currency_id),
		Color:        color,
		IconCode:     database.TypeIcons(icon_id),
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

func ChangeAccountColorHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
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

	color := r.FormValue("Color")

	err = database.ChangeAccountColor(db, accountID, color)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to change account's color: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("Color in account was changed successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ChangeAccountIconHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
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

	iconID, err := strconv.Atoi(r.FormValue("IconID"))
	if err != nil {
		http.Error(w, "problem with icon_id. use normal values", http.StatusBadRequest)
		return
	}

	err = database.ChangeAccountIcon(db, accountID, iconID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to change account's icon: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Icon for account %d was changed successfully to %d", accountID, iconID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
