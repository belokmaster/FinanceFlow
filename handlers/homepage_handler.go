package handlers

import (
	"net/http"
	"text/template"

	"finance_flow/database"

	"gorm.io/gorm"
)

func HomeHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB, path string) {
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	accounts, err := database.GetAccounts(db)
	if err != nil {
		http.Error(w, "could not get accounts", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, "could not parse template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, accounts)
	if err != nil {
		http.Error(w, "could not execute template", http.StatusInternalServerError)
	}
}
