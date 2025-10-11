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
	ID   uint   `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"uniqueIndex;not null"`
}

type SubCategory struct {
	ID         uint   `gorm:"primaryKey;autoIncrement"`
	Name       string `gorm:"uniqueIndex;not null"`
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
	SubCategoryID uint
	SubCategory   SubCategory
	Type          TypeTransaction
	Amount        float64
	Comment       string
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

type TypeIcons int

const (
	Coin TypeIcons = iota
	Mark
	Card
	Wallet
	House
	Bag
	Bonus
)

var IconFiles = map[TypeIcons]string{
	Coin:   "Наличные",
	Mark:   "Сберегательный счет",
	Card:   "Карта",
	Wallet: "Общий счет",
	House:  "Ипотека",
	Bag:    "Заем",
	Bonus:  "Бонусы",
}

// kostil...
var IconNamesToIDs = map[string]TypeIcons{
	"Наличные":            Coin,
	"Сберегательный счет": Mark,
	"Карта":               Card,
	"Общий счет":          Wallet,
	"Ипотека":             House,
	"Заем":                Bag,
	"Бонусы":              Bonus,
}
