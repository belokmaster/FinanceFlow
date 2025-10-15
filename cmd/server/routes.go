package main

import (
	"finance_flow/internal/handlers"
	"net/http"

	"gorm.io/gorm"
)

func setupRoutes(db *gorm.DB) {
	// static routes
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// GET routes
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

	// POST routes
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

	http.HandleFunc("/delete_transaction", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteTransactionHandler(w, r, db)
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

	http.HandleFunc("/update_transaction", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateTransactionHandler(w, r, db)
	})
}
