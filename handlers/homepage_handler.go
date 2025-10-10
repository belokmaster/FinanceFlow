package handlers

import (
	"log"
	"net/http"
	"text/template"

	"finance_flow/database"
	"finance_flow/icons"

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

	var accountsForView []AccountView
	for _, acc := range accounts {
		symbol, ok := database.CurrencySymbols[acc.CurrencyCode]
		if !ok {
			symbol = "?"
		}

		color := "#" + acc.Color
		if color == "" {
			color = "#b1afafff"
		}

		iconFileName, ok := database.IconFiles[acc.IconCode]
		if !ok {
			iconFileName = "default"
		}

		iconHTML, ok := icons.IconCache[iconFileName]
		if !ok {
			iconHTML = icons.IconCache["coin"]
		}

		accountsForView = append(accountsForView, AccountView{
			Name:           acc.Name,
			Balance:        acc.Balance,
			CurrencySymbol: symbol,
			Color:          color,
			Icon:           iconHTML,
		})
	}

	tmpl, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, "could not parse template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, accountsForView)
	if err != nil {
		log.Printf("error while proccess working: %v", err)
	}
}
