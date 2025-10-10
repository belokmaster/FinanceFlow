package handlers

import (
	"html/template"
)

type AccountView struct {
	Name           string
	Balance        float64
	CurrencySymbol string
	Color          string
	Icon           template.HTML
}
