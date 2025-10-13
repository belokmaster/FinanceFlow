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

var IconAccountFiles = map[TypeIcons]string{
	Coin:   "Наличные",
	Mark:   "Сберегательный счет",
	Card:   "Карта",
	Wallet: "Общий счет",
	House:  "Ипотека",
	Bag:    "Заем",
	Bonus:  "Бонусы",
}

// kostil...
var IconAccountNamesToIDs = map[string]TypeIcons{
	"Наличные":            Coin,
	"Сберегательный счет": Mark,
	"Карта":               Card,
	"Общий счет":          Wallet,
	"Ипотека":             House,
	"Заем":                Bag,
	"Бонусы":              Bonus,
}

type TypeCategoryIcons int

const (
	Food TypeCategoryIcons = iota
	Housing1
	Housing2
	Car_Repair1
	Car_Repair2
	Coin_cat
	Computer1
	Computer2
	Bus1
	Bus2
	Community1
	Community2
	Community3
	Expense1
	Expense2
)

var IconCategoryFiles = map[TypeCategoryIcons]string{
	Food:        "Еда",
	Housing1:    "Жилье1",
	Housing2:    "Жилье2",
	Car_Repair1: "Транспортное средство1",
	Car_Repair2: "Транспортное средство2",
	Coin_cat:    "Доход",
	Computer1:   "Связь, ПК1",
	Computer2:   "Связь, ПК2",
	Bus1:        "Общественный транспорт1",
	Bus2:        "Общественный транспорт2",
	Community1:  "Жизнь и развлечения1",
	Community2:  "Жизнь и развлечения2",
	Community3:  "Жизнь и развлечения3",
	Expense1:    "Финансовые расходы1",
	Expense2:    "Финансовые расходы2",
}

var IconCategoryNamesToIDs = map[string]TypeCategoryIcons{
	"Еда":    Food,
	"Жилье1": Housing1,
	"Жилье2": Housing2,
	"Транспортное средство1": Car_Repair1,
	"Транспортное средство2": Car_Repair2,
	"Доход":      Coin_cat,
	"Связь, ПК1": Computer1,
	"Связь, ПК2": Computer2,
	"Общественный транспорт1": Bus1,
	"Общественный транспорт2": Bus2,
	"Жизнь и развлечения1":    Community1,
	"Жизнь и развлечения2":    Community2,
	"Жизнь и развлечения3":    Community3,
	"Финансовые расходы1":     Expense1,
	"Финансовые расходы2":     Expense2,
}

type TypeSubCategoryIcons int

const (
	Restaurant TypeSubCategoryIcons = iota
)

var IconSubCategoryFiles = map[TypeSubCategoryIcons]string{
	Restaurant: "Ресторан1",
}

var IconSubCategoryNamesToIDs = map[string]TypeSubCategoryIcons{
	"Ресторан1": Restaurant,
}
