package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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

		accountID, _ := strconv.Atoi(r.FormValue("AccountID"))
		categoryID, _ := strconv.Atoi(r.FormValue("CategoryID"))
		subCategoryID, _ := strconv.Atoi(r.FormValue("SubCategoryID"))
		transType, _ := strconv.Atoi(r.FormValue("Type"))
		amount, _ := strconv.ParseFloat(r.FormValue("Amount"), 64)
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
