package icons

import "html/template"

var accountFileToName = map[string]string{
	"wallet": "Общий счет",
	"card":   "Карта",
	"coin":   "Наличные",
	"mark":   "Сберегательный счет",
	"house":  "Ипотека",
	"bag":    "Заем",
	"bonus":  "Бонусы",
}

var categoryFileToName = map[string]string{
	"food":        "Еда",
	"housing1":    "Жилье1",
	"housing2":    "Жилье2",
	"car_repair1": "Транспортное средство1",
	"car_repair2": "Транспортное средство2",
	"coin_cat":    "Доход",
	"computer1":   "Связь, ПК1",
	"computer2":   "Связь, ПК2",
	"bus1":        "Общественный транспорт1",
	"bus2":        "Общественный транспорт2",
	"community1":  "Жизнь и развлечения1",
	"community2":  "Жизнь и развлечения2",
	"community3":  "Жизнь и развлечения3",
	"expense1":    "Финансовые расходы1",
	"expense2":    "Финансовые расходы2",
}

var subCategoryFileToName = map[string]string{
	"restaurant1": "Ресторан1",
	"pizza1":      "Фаст-фуд1",
}

var AccountIconCache = make(map[string]template.HTML)
var CategoryIconCache = make(map[string]template.HTML)
var SubCategoryIconCache = make(map[string]template.HTML)
