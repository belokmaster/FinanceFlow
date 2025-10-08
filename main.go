package main

import (
	"log"
	"net/http"
	"os"

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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/create_account", func(w http.ResponseWriter, r *http.Request) {
		CreateAccountHandler(w, r, db)
	})

	http.HandleFunc("/delete_account", func(w http.ResponseWriter, r *http.Request) {
		DeleteAccountHandler(w, r, db)
	})

	http.HandleFunc("/transfer", func(w http.ResponseWriter, r *http.Request) {
		TransferHandler(w, r, db)
	})

	http.HandleFunc("/submit_transaction", func(w http.ResponseWriter, r *http.Request) {
		TransactionHandler(w, r, db)
	})

	http.HandleFunc("/add_category", func(w http.ResponseWriter, r *http.Request) {
		CreateCategoryHandler(w, r, db)
	})

	http.HandleFunc("/delete_category", func(w http.ResponseWriter, r *http.Request) {
		DeleteCategoryHandler(w, r, db)
	})

	http.HandleFunc("/add_subcategory", func(w http.ResponseWriter, r *http.Request) {
		CreateSubCategoryHandler(w, r, db)
	})

	http.HandleFunc("/delete_subcategory", func(w http.ResponseWriter, r *http.Request) {
		DeleteSubCategoryHandler(w, r, db)
	})

	log.Println("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
