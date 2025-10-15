package handlers

import (
	"finance_flow/internal/database"
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
	ID            uint
	Name          string
	Color         string
	IconKey       string
	IconHTML      template.HTML
	Subcategories []SubCategoryView
}

type CategoryPageData struct {
	Categories       []CategoryView
	CategoryIcons    map[string]template.HTML
	SubcategoryIcons map[string]template.HTML
}

type SubCategoryView struct {
	ID       uint
	Name     string
	Color    string
	IconKey  string
	IconHTML template.HTML
	ParentID uint
}

type SubCategoryPageData struct {
	Categories []SubCategoryView
	Icons      map[string]template.HTML
}

type GroupedTransactions struct {
	Date           string
	TotalAmount    float64
	CurrencySymbol string
	Transactions   []TransactionView
}

type TransactionPageData struct {
	Accounts            []AccountView
	Categories          []CategoryView
	Transactions        []TransactionView
	GroupedTransactions []GroupedTransactions
}

type TransactionView struct {
	ID                 uint
	Type               database.TypeTransaction
	Amount             float64
	AccountID          uint
	AccountName        string
	AccountColor       string
	CurrencySymbol     string
	CategoryID         uint
	CategoryName       string
	CategoryColor      string
	CategoryIconHTML   template.HTML
	DisplayName        string
	ParentCategoryName string
	SubCategoryID      *uint
	SubCategoryName    *string
	Date               string
	FormattedTime      string
	FormattedDate      string
	Description        string
}
