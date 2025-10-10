package database

import "time"

type Account struct {
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	Name         string `gorm:"uniqueIndex;not null"`
	Balance      float64
	CurrencyCode TypeCurrency
	Color        string
	IconCode     TypeIcons
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
)

var CurrencySymbols = map[TypeCurrency]string{
	Ruble:  "₽",
	Dollar: "$",
	Euro:   "€",
	Yuan:   "¥",
}

type TypeIcons int

const (
	Nothing TypeIcons = iota
	Wallet
	Card
	Cash
	Education
	Shopping
	Mark
	Invest
	Money
	Case
	Coin
	Gift
	Heart
	Diamond
	Plane
	Hova
)

var IconSymbols = map[TypeIcons]string{
	Nothing:   "",
	Wallet:    "💰",
	Card:      "💳",
	Cash:      "💵",
	Education: "🎓",
	Shopping:  "🛍️",
	Mark:      "🐷",
	Invest:    "📈",
	Money:     "💸",
	Case:      "💼",
	Coin:      "🪙",
	Gift:      "🎁",
	Heart:     "❤️",
	Diamond:   "💎",
	Plane:     "✈️",
	Hova:      "🤿",
}
