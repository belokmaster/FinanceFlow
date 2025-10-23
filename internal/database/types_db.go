package database

import "time"

type Account struct {
	ID           uint         `gorm:"primaryKey;autoIncrement"`
	Name         string       `gorm:"uniqueIndex;not null"`
	Balance      float64      `gorm:"type:decimal(15,2)"`
	CurrencyCode TypeCurrency `gorm:"column:currency_code"`
	Color        string       `gorm:"column:color"`
	IconCode     TypeIcons    `gorm:"column:icon_code"`
}

type Category struct {
	ID            uint              `gorm:"primaryKey;autoIncrement"`
	Name          string            `gorm:"uniqueIndex;not null"`
	Color         string            `gorm:"column:color"`
	IconCode      TypeCategoryIcons `gorm:"column:icon_code"`
	SubCategories []SubCategory     `gorm:"foreignKey:CategoryID"`
}

type SubCategory struct {
	ID         uint                 `gorm:"primaryKey;autoIncrement"`
	Name       string               `gorm:"uniqueIndex;not null"`
	Color      string               `gorm:"column:color"`
	IconCode   TypeSubCategoryIcons `gorm:"column:icon_code"`
	CategoryID uint
	Category   Category `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type TypeTransaction int

const (
	Income TypeTransaction = iota
	Expense
	Transfer
)

type Transaction struct {
	ID            uint `gorm:"primaryKey;autoIncrement"`
	AccountID     uint
	Account       Account
	CategoryID    uint
	Category      Category
	SubCategoryID uint        `gorm:"null"`
	SubCategory   SubCategory `gorm:"foreignKey:SubCategoryID"`
	Type          TypeTransaction
	Amount        float64
	Comment       string `gorm:"size:500;null"`
	Date          time.Time
}

type TransferTransaction struct {
	ID                uint `gorm:"primaryKey;autoIncrement"`
	AccountID         uint
	Account           Account
	TransferAccountID uint
	TransferAccount   Account
	Type              TypeTransaction
	Amount            float64
	Comment           string
	Date              time.Time
}

type TypeCurrency int

const (
	Ruble TypeCurrency = iota
	Dollar
	Euro
	Yuan
	GBPound
	Rupee
)

var CurrencySymbols = map[TypeCurrency]string{
	Ruble:   "₽",
	Dollar:  "$",
	Euro:    "€",
	Yuan:    "¥",
	GBPound: "£",
	Rupee:   "₹",
}
