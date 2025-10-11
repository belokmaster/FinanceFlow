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

type CategoryView struct {
	ID       uint
	Name     string
	Color    string
	IconKey  string
	IconHTML template.HTML
}

type CategoryPageData struct {
	Categories []CategoryView
	Icons      map[string]template.HTML
}
