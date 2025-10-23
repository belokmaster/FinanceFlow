package database

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
	Restaurant1 TypeSubCategoryIcons = iota
	FastFood1
	Cooking1
	Cooking2
	Pharmacy1
	Baby1
	Pets1
)

var IconSubCategoryFiles = map[TypeSubCategoryIcons]string{
	Restaurant1: "Ресторан1",
	FastFood1:   "Фаст-фуд1",
	Cooking1:    "Готовка1",
	Cooking2:    "Десерт1",
	Pharmacy1:   "Аптека1",
	Baby1:       "Детские товары1",
	Pets1:       "Домашние животные, питомцы1",
}

var IconSubCategoryNamesToIDs = map[string]TypeSubCategoryIcons{
	"Ресторан1":       Restaurant1,
	"Фаст-фуд1":       FastFood1,
	"Готовка1":        Cooking1,
	"Десерт1":         Cooking2,
	"Аптека1":         Pharmacy1,
	"Детские товары1": Baby1,
	"Домашние животные, питомцы1": Pets1,
}
