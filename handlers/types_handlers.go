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
	Icon           template.HTML
}
