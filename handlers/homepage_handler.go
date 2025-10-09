package handlers

import (
	"net/http"
	"text/template"

	"finance_flow/database"

	"gorm.io/gorm"
)

type AccountView struct {
	Name           string
	Balance        float64
	CurrencySymbol string
}

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

	var accountsForView []AccountView
	for _, acc := range accounts {
		symbol, ok := database.CurrencySymbols[acc.CurrencyCode]
		if !ok {
			symbol = "?"
		}

		accountsForView = append(accountsForView, AccountView{
			Name:           acc.Name,
			Balance:        acc.Balance,
			CurrencySymbol: symbol,
		})
	}

	tmpl, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, "could not parse template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, accountsForView)
	if err != nil {
		http.Error(w, "could not execute template", http.StatusInternalServerError)
	}
}
