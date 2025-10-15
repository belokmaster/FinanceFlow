package handlers

import (
	"finance_flow/internal/database"
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

	name := r.FormValue("Name")
	name = strings.TrimSpace(name) // why it works??

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

	if currency_id >= len(database.CurrencySymbols) {
		http.Error(w, "invalid currency id", http.StatusBadRequest)
		return
	}

	color := r.FormValue("Color")
	log.Printf("CreateAccountHandler: Color '%s'", color)
	color = strings.TrimPrefix(color, "#")

	if color == "" {
		color = "4cd67a"
	}

	icon_id, err := strconv.Atoi(r.FormValue("Icon"))
	if err != nil {
		log.Printf("CreateAccountHandler: Invalid icon ID '%s': %v", r.FormValue("Icon"), err)
		http.Error(w, "problem with icon. use normal values", http.StatusBadRequest)
		return
	}

	if icon_id < 0 {
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

	log.Printf("CreateAccountHandler: Creating account with Name=%s, Balance=%.2f, Currency=%d, Icon=%d", name, balance, currency_id, icon_id)

	err = database.CreateNewAccount(db, acc)
	if err != nil {
		log.Printf("CreateAccountHandler: Failed to create account: %v", err)
		http.Error(w, fmt.Sprintf("failed to create a new acc: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("CreateAccountHandler: Account created successfully - Name=%s", name)
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

func UpdateAccountHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.Method != http.MethodPost {
		log.Printf("UpdateAccountHandler: Invalid method %s, redirecting to /", r.Method)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Printf("UpdateAccountHandler: Error parsing form: %v", err)
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	log.Printf("UpdateAccountHandler: Form data received - %+v", r.Form)

	accountID, err := strconv.Atoi(r.FormValue("ID"))
	if err != nil {
		log.Printf("UpdateAccountHandler: Invalid account ID '%s': %v", r.FormValue("ID"), err)
		http.Error(w, "problem with id. use normal values", http.StatusBadRequest)
		return
	}

	newColor := r.FormValue("Color")
	if newColor != "" {
		err = database.ChangeAccountColor(db, accountID, newColor)
		if err != nil {
			log.Printf("UpdateAccountHandler: Changing color to '%s' for account %d", newColor, accountID)
			http.Error(w, fmt.Sprintf("failed to change account's color: %v", err), http.StatusInternalServerError)
			return
		}
	}

	newName := r.FormValue("Name")
	if newName != "" {
		err = database.ChangeAccountName(db, accountID, newName)
		if err != nil {
			log.Printf("UpdateAccountHandler: Changing name to '%s' for account %d", newName, accountID)
			http.Error(w, fmt.Sprintf("failed to change account's name: %v", err), http.StatusInternalServerError)
			return
		}
	}

	icon := r.FormValue("Icon")
	if icon != "" {
		iconCode, exists := database.IconAccountNamesToIDs[icon]
		if !exists {
			log.Printf("UpdateAccountHandler: Unknown icon name '%s'", icon)
			http.Error(w, "unknown icon name", http.StatusBadRequest)
			return
		}

		log.Printf("UpdateAccountHandler: Changing icon for account ID %d to icon ID %d", accountID, int(iconCode))

		err = database.ChangeAccountIcon(db, accountID, int(iconCode))
		if err != nil {
			log.Printf("UpdateAccountHandler: Failed to change icon for account ID %d: %v", accountID, err)
			http.Error(w, fmt.Sprintf("failed to change account's icon: %v", err), http.StatusInternalServerError)
			return
		}
	}

	newBalance, err := strconv.ParseFloat(r.FormValue("Balance"), 64)
	if err != nil {
		log.Printf("UpdateAccountHandler: Invalid amount '%s': %v", r.FormValue("Balance"), err)
		http.Error(w, "problem with balance. use normal values", http.StatusBadRequest)
		return
	}

	err = database.ChangeAccountBalance(db, accountID, newBalance)
	if err != nil {
		log.Printf("UpdateAccountHandler: Failed to change balance for account ID %d: %v", accountID, err)
		http.Error(w, fmt.Sprintf("failed to change account's balance: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("UpdateAccountHandler: Successfully processed all changes for account %d", accountID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
