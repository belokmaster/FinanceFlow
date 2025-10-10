package handlers

import (
	"html/template"
)

type AccountView struct {
	ID             uint
	Name           string
	Balance        float64
	CurrencySymbol string
	Color          string
	IconKey        string
	IconHTML       template.HTML
}

type HomePageData struct {
	Accounts []AccountView
	Icons    map[string]template.HTML
}
