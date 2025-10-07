package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

func main() {
	dbPath := "database.db" // mb change this point later

	// check if db exists
	_, err := os.Stat(dbPath)
	dbExists := !os.IsNotExist(err)

	// open the db using gorm
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("problem connecting to db: %v", err)
	}

	// initialize a new db if it doesnt exist
	if !dbExists {
		log.Println("initialize a new bd")
		err := InitDB(db)
		if err != nil {
			log.Fatalf("failed to initialize the bd: %v", err)
		}
	}

	// close the db
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("problem to close to db: %v", err)
	}
	defer sqlDB.Close()

	log.Println("The database opens successfully")

	// main page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
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
	})

	http.HandleFunc("/transfer", func(w http.ResponseWriter, r *http.Request) {
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
	})

	// add trans in server
	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
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
	})

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
