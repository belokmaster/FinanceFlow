package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

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

	name := r.FormValue("Name")
	name = strings.TrimSpace(name) // why it works????????

	balance, err := strconv.ParseFloat(r.FormValue("Balance"), 64)
	if err != nil {
		http.Error(w, "problem with balance. use normal values", http.StatusBadRequest)
		return
	}

	acc := Account{
		ID:      uint(accountID),
		Name:    name,
		Balance: balance,
	}

	err = CreateNewAccount(db, acc)
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

	err = DeleteAccount(db, accountID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to delete a acc: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("A account was deleted sucessfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

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

	transferAccountID, err := strconv.Atoi(r.FormValue("TransferAccountID"))
	if err != nil {
		http.Error(w, "problem with transferaccount_id. use normal values", http.StatusBadRequest)
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

	tx := TransferTransaction{
		AccountID:         uint(accountID),
		TransferAccountID: uint(transferAccountID),
		Type:              TypeTransaction(transType),
		Amount:            amount,
		Comment:           comment,
		Date:              time.Now(),
	}

	err = TransferToAnotherAcc(db, tx)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to transfer: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("Transfer added and accounts balance updated")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

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

	categoryID, err := strconv.Atoi(r.FormValue("CategoryID"))
	if err != nil {
		http.Error(w, "problem with category_id. use normal values", http.StatusBadRequest)
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

	tx := Transaction{
		AccountID:     uint(accountID),
		CategoryID:    uint(categoryID),
		SubCategoryID: uint(subCategoryID),
		Type:          TypeTransaction(transType),
		Amount:        amount,
		Comment:       comment,
		Date:          time.Now(),
	}

	err = AddTransaction(db, tx)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to add transaction: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("Transaction added and account balance updated")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func CreateCategoryHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	categoryID, err := strconv.Atoi(r.FormValue("ID"))
	if err != nil {
		http.Error(w, "problem with id. use normal values", http.StatusBadRequest)
		return
	}

	if categoryID < 0 {
		http.Error(w, "problem with id. use positive values", http.StatusBadRequest)
		return
	}

	name := r.FormValue("Name")
	name = strings.TrimSpace(name)

	category := Category{
		ID:   uint(categoryID),
		Name: name,
	}

	err = AddCategory(db, category)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create a new category: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("A new category was created sucessfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteCategoryHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	categoryID, err := strconv.Atoi(r.FormValue("ID"))
	if err != nil {
		http.Error(w, "problem with id. use normal values", http.StatusBadRequest)
		return
	}

	err = DeleteCategory(db, categoryID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to delete category: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("category was deleted sucessfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func CreateSubCategoryHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	subCategoryID, err := strconv.Atoi(r.FormValue("ID"))
	if err != nil {
		http.Error(w, "problem with id. use normal values", http.StatusBadRequest)
		return
	}

	if subCategoryID < 0 {
		http.Error(w, "problem with id. use positive values", http.StatusBadRequest)
		return
	}

	categoryID, err := strconv.Atoi(r.FormValue("ID"))
	if err != nil {
		http.Error(w, "problem with id. use normal values", http.StatusBadRequest)
		return
	}

	if categoryID < 0 {
		http.Error(w, "problem with id. use positive values", http.StatusBadRequest)
		return
	}

	name := r.FormValue("Name")
	name = strings.TrimSpace(name)

	subCategory := SubCategory{
		ID:         uint(subCategoryID),
		Name:       name,
		CategoryID: uint(categoryID),
	}

	err = AddSubCategory(db, subCategory)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create a new category: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("A new sub_category was created sucessfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteSubCategoryHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	subCategoryID, err := strconv.Atoi(r.FormValue("ID"))
	if err != nil {
		http.Error(w, "problem with id. use normal values", http.StatusBadRequest)
		return
	}

	err = DeleteSubCategory(db, subCategoryID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to delete sub_category: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("sub_category was deleted sucessfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
