package main

import (
	"finance_flow/internal/database"
	"finance_flow/internal/icons"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

func main() {
	db := database.InitDatabase()
	defer database.CloseDatabase(db)
	icons.InitIcons()

	setupRoutes(db)

	log.Println("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
