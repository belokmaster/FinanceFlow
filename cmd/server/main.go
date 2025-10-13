package main

import (
	"finance_flow/internal/database"
	"finance_flow/internal/handlers"
	"finance_flow/internal/icons"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

func main() {
	dbPath := "internal/database/database.db" // mb change this point later

	// check if db exists
	_, err := os.Stat(dbPath)
	dbExists := !os.IsNotExist(err)

	// open the db using gorm
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("problem connecting to db: %v", err)
	}

	icons.InitIcons()

	// initialize a new db if it doesnt exist
	if !dbExists {
		log.Println("initialize a new bd")
		err := database.InitDB(db)
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

	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := "web/templates/index.html"
		handlers.HomeHandler(w, r, db, path)
	})

	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		path := "web/templates/categories.html"
		handlers.CategoryPageHandler(w, r, db, path)
	})

	http.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		path := "web/templates/transaction.html"
		handlers.NewTransactionPageHandler(w, r, db, path)
	})

	http.HandleFunc("/create_account", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateAccountHandler(w, r, db)
	})

	http.HandleFunc("/delete_account", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteAccountHandler(w, r, db)
	})

	http.HandleFunc("/transfer", func(w http.ResponseWriter, r *http.Request) {
		handlers.TransferHandler(w, r, db)
	})

	http.HandleFunc("/submit_transaction", func(w http.ResponseWriter, r *http.Request) {
		handlers.TransactionHandler(w, r, db)
	})

	http.HandleFunc("/create_category", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateCategoryHandler(w, r, db)
	})

	http.HandleFunc("/delete_category", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteCategoryHandler(w, r, db)
	})

	http.HandleFunc("/create_subcategory", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateSubCategoryHandler(w, r, db)
	})

	http.HandleFunc("/delete_subcategory", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteSubCategoryHandler(w, r, db)
	})

	http.HandleFunc("/update_account", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateAccountHandler(w, r, db)
	})

	http.HandleFunc("/update_category", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateCategoryHandler(w, r, db)
	})

	http.HandleFunc("/update_subcategory", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateSubCategoryHandler(w, r, db)
	})

	log.Println("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
