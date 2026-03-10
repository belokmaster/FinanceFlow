package handlers

import (
	"finance_flow/internal/database"
	"fmt"
	"log"
	"net/http"

	"gorm.io/gorm"
)

func GenerateTestDataHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := database.GenerateTestData(db); err != nil {
		log.Printf("GenerateTestDataHandler: failed to generate demo data: %v", err)
		http.Error(w, fmt.Sprintf("failed to generate demo data: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("GenerateTestDataHandler: demo data generated successfully")
	http.Redirect(w, r, "/transactions", http.StatusSeeOther)
}
