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
	log.Printf("CreateAccountHandler: Starting account creation")

	if r.Method != http.MethodPost {
		log.Printf("CreateAccountHandler: Invalid method %s, redirecting to /", r.Method)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Printf("CreateAccountHandler: Error parsing form: %v", err)
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	log.Printf("CreateAccountHandler: Form data received - %+v", r.Form)

	accountID, err := strconv.Atoi(r.FormValue("ID"))
	if err != nil {
		log.Printf("CreateAccountHandler: Invalid account ID '%s': %v", r.FormValue("ID"), err)
		http.Error(w, "problem with id. use normal values", http.StatusBadRequest)
		return
	}

	if accountID < 0 {
		log.Printf("CreateAccountHandler: Negative account ID %d", accountID)
		http.Error(w, "problem with id. use positive values", http.StatusBadRequest)
		return
	}

	name := r.FormValue("Name")
	name = strings.TrimSpace(name) // why it works????????
	log.Printf("CreateAccountHandler: Account name '%s' (trimmed)", name)

	balance, err := strconv.ParseFloat(r.FormValue("Balance"), 64)
	if err != nil {
		log.Printf("CreateAccountHandler: Invalid balance '%s': %v", r.FormValue("Balance"), err)
		http.Error(w, "problem with balance. use normal values", http.StatusBadRequest)
		return
	}

	currency_id, err := strconv.Atoi(r.FormValue("Currency"))
	if err != nil {
		log.Printf("CreateAccountHandler: Invalid currency '%s': %v", r.FormValue("Currency"), err)
		http.Error(w, "problem with currency. use normal values", http.StatusBadRequest)
		return
	}

	if currency_id < 0 {
		log.Printf("CreateAccountHandler: Negative currency ID %d", currency_id)
		http.Error(w, "problem with currency. use positive values", http.StatusBadRequest)
		return
	}

	color := r.FormValue("Color")
	log.Printf("CreateAccountHandler: Color '%s'", color)

	icon_id, err := strconv.Atoi(r.FormValue("Icon"))
	if err != nil {
		log.Printf("CreateAccountHandler: Invalid icon ID '%s': %v", r.FormValue("Icon"), err)
		http.Error(w, "problem with icon. use normal values", http.StatusBadRequest)
		return
	}

	if currency_id < 0 {
		log.Printf("CreateAccountHandler: Negative icon ID %d", icon_id)
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

	log.Printf("CreateAccountHandler: Creating account with ID=%d, Name=%s, Balance=%.2f, Currency=%d, Icon=%d", accountID, name, balance, currency_id, icon_id)

	err = database.CreateNewAccount(db, acc)
	if err != nil {
		log.Printf("CreateAccountHandler: Failed to create account: %v", err)
		http.Error(w, fmt.Sprintf("failed to create a new acc: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("CreateAccountHandler: Account created successfully - ID=%d, Name=%s", accountID, name)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteAccountHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	log.Printf("DeleteAccountHandler: Starting account deleting")

	if r.Method != http.MethodPost {
		log.Printf("DeleteAccountHandler: Invalid method %s, redirecting to /", r.Method)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Printf("DeleteAccountHandler: Error parsing form: %v", err)
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	log.Printf("DeleteAccountHandler: Form data received - %+v", r.Form)

	accountID, err := strconv.Atoi(r.FormValue("ID"))
	if err != nil {
		log.Printf("DeleteAccountHandler: Invalid account ID '%s': %v", r.FormValue("ID"), err)
		http.Error(w, "problem with id. use normal values", http.StatusBadRequest)
		return
	}

	if accountID < 0 {
		log.Printf("DeleteAccountHandler: Negative account ID %d", accountID)
		http.Error(w, "problem with id. use positive values", http.StatusBadRequest)
		return
	}

	log.Printf("DeleteAccountHandler: Attempting to delete account ID %d", accountID)

	err = database.DeleteAccount(db, accountID)
	if err != nil {
		log.Printf("DeleteAccountHandler: Failed to delete account ID %d: %v", accountID, err)
		http.Error(w, fmt.Sprintf("failed to delete a acc: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("DeleteAccountHandler: Account ID %d deleted successfully", accountID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ChangeAccountColorHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	log.Printf("ChangeAccountColorHandler: Starting color change")

	if r.Method != http.MethodPost {
		log.Printf("ChangeAccountColorHandler: Invalid method %s, redirecting to /", r.Method)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Printf("ChangeAccountColorHandler: Error parsing form: %v", err)
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	log.Printf("ChangeAccountColorHandler: Form data received - %+v", r.Form)

	accountID, err := strconv.Atoi(r.FormValue("ID"))
	if err != nil {
		log.Printf("ChangeAccountColorHandler: Invalid account ID '%s': %v", r.FormValue("ID"), err)
		http.Error(w, "problem with id. use normal values", http.StatusBadRequest)
		return
	}

	color := r.FormValue("Color")
	log.Printf("ChangeAccountColorHandler: Changing color for account ID %d to '%s'", accountID, color)

	err = database.ChangeAccountColor(db, accountID, color)
	if err != nil {
		log.Printf("ChangeAccountColorHandler: Failed to change color for account ID %d: %v", accountID, err)
		http.Error(w, fmt.Sprintf("failed to change account's color: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("ChangeAccountColorHandler: Color changed successfully for account ID %d", accountID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ChangeAccountIconHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	log.Printf("ChangeAccountIconHandler: Starting icon change")

	if r.Method != http.MethodPost {
		log.Printf("ChangeAccountIconHandler: Invalid method %s, redirecting to /", r.Method)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Printf("ChangeAccountIconHandler: Error parsing form: %v", err)
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	log.Printf("ChangeAccountIconHandler: Form data received - %+v", r.Form)

	accountID, err := strconv.Atoi(r.FormValue("ID"))
	if err != nil {
		log.Printf("ChangeAccountIconHandler: Invalid account ID '%s': %v", r.FormValue("ID"), err)
		http.Error(w, "problem with id. use normal values", http.StatusBadRequest)
		return
	}

	iconID, err := strconv.Atoi(r.FormValue("IconID"))
	if err != nil {
		log.Printf("ChangeAccountIconHandler: Invalid icon ID '%s': %v", r.FormValue("IconID"), err)
		http.Error(w, "problem with icon_id. use normal values", http.StatusBadRequest)
		return
	}

	log.Printf("ChangeAccountIconHandler: Changing icon for account ID %d to icon ID %d", accountID, iconID)

	err = database.ChangeAccountIcon(db, accountID, iconID)
	if err != nil {
		log.Printf("ChangeAccountIconHandler: Failed to change icon for account ID %d: %v", accountID, err)
		http.Error(w, fmt.Sprintf("failed to change account's icon: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("ChangeAccountIconHandler: Icon changed successfully for account ID %d to icon ID %d", accountID, iconID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
