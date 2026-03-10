package database

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

func GenerateTestData(db *gorm.DB) error {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	type accountSeed struct {
		key      string
		name     string
		balance  float64
		currency TypeCurrency
		color    string
		icon     TypeIcons
	}

	accountSeeds := []accountSeed{
		{key: "cash", name: "Наличные", balance: 30000, currency: Ruble, color: "4cd67a", icon: Coin},
		{key: "card", name: "Основная карта", balance: 180000, currency: Ruble, color: "4aa8ff", icon: Card},
		{key: "savings", name: "Сбережения", balance: 420000, currency: Ruble, color: "ffb84d", icon: Mark},
	}

	accountsByKey := make(map[string]Account, len(accountSeeds))
	for _, seed := range accountSeeds {
		name, err := ensureUniqueName(db, &Account{}, seed.name)
		if err != nil {
			return fmt.Errorf("failed to generate unique account name: %w", err)
		}

		acc := Account{
			Name:         name,
			Balance:      seed.balance,
			CurrencyCode: seed.currency,
			Color:        seed.color,
			IconCode:     seed.icon,
		}
		if err := CreateNewAccount(db, acc); err != nil {
			return fmt.Errorf("failed to create account %s: %w", name, err)
		}

		var created Account
		if err := db.Where("name = ?", name).First(&created).Error; err != nil {
			return fmt.Errorf("failed to fetch created account %s: %w", name, err)
		}
		accountsByKey[seed.key] = created
	}

	type categorySeed struct {
		key   string
		name  string
		color string
		icon  TypeCategoryIcons
	}

	categorySeeds := []categorySeed{
		{key: "income", name: "Доход", color: "4cd67a", icon: Coin_cat},
		{key: "food", name: "Еда", color: "ff7f50", icon: Food},
		{key: "transport", name: "Транспорт", color: "5f9ea0", icon: Bus1},
		{key: "home", name: "Дом", color: "7f8fa6", icon: Housing1},
		{key: "fun", name: "Развлечения", color: "ff6b81", icon: Community2},
		{key: "health", name: "Здоровье", color: "2ecc71", icon: Expense1},
		{key: "shopping", name: "Покупки", color: "9b59b6", icon: Community1},
	}

	categoriesByKey := make(map[string]Category, len(categorySeeds))
	for _, seed := range categorySeeds {
		name, err := ensureUniqueName(db, &Category{}, seed.name)
		if err != nil {
			return fmt.Errorf("failed to generate unique category name: %w", err)
		}

		cat := Category{Name: name, Color: seed.color, IconCode: seed.icon}
		if err := CreateCategory(db, cat); err != nil {
			return fmt.Errorf("failed to create category %s: %w", name, err)
		}

		var created Category
		if err := db.Where("name = ?", name).First(&created).Error; err != nil {
			return fmt.Errorf("failed to fetch created category %s: %w", name, err)
		}
		categoriesByKey[seed.key] = created
	}

	type subcategorySeed struct {
		key      string
		parent   string
		name     string
		color    string
		iconCode TypeSubCategoryIcons
	}

	subcategorySeeds := []subcategorySeed{
		{key: "groceries", parent: "food", name: "Продукты", color: "ff9f43", iconCode: Cooking1},
		{key: "cafe", parent: "food", name: "Кафе и рестораны", color: "ee5253", iconCode: Restaurant1},
		{key: "fastfood", parent: "food", name: "Фастфуд", color: "f39c12", iconCode: FastFood1},
		{key: "taxi", parent: "transport", name: "Такси", color: "48c9b0", iconCode: Restaurant1},
		{key: "public", parent: "transport", name: "Общественный транспорт", color: "3498db", iconCode: Pets1},
		{key: "pharmacy", parent: "health", name: "Аптека", color: "58d68d", iconCode: Pharmacy1},
		{key: "marketplace", parent: "shopping", name: "Маркетплейсы", color: "a569bd", iconCode: Baby1},
	}

	subcategoriesByKey := make(map[string]SubCategory, len(subcategorySeeds))
	for _, seed := range subcategorySeeds {
		parent, ok := categoriesByKey[seed.parent]
		if !ok {
			return fmt.Errorf("parent category with key %s not found", seed.parent)
		}

		name, err := ensureUniqueName(db, &SubCategory{}, seed.name)
		if err != nil {
			return fmt.Errorf("failed to generate unique subcategory name: %w", err)
		}

		sub := SubCategory{
			Name:       name,
			Color:      seed.color,
			IconCode:   seed.iconCode,
			CategoryID: parent.ID,
		}
		if err := AddSubCategory(db, sub); err != nil {
			return fmt.Errorf("failed to create subcategory %s: %w", name, err)
		}

		var created SubCategory
		if err := db.Where("name = ?", name).First(&created).Error; err != nil {
			return fmt.Errorf("failed to fetch created subcategory %s: %w", name, err)
		}
		subcategoriesByKey[seed.key] = created
	}

	type txPattern struct {
		categoryKey    string
		subcategoryKey string
		accountKey     string
		minAmount      float64
		maxAmount      float64
		comments       []string
		weight         int
	}

	expensePatterns := []txPattern{
		{categoryKey: "food", subcategoryKey: "groceries", accountKey: "card", minAmount: 350, maxAmount: 4200, comments: []string{"Продукты", "Супермаркет", "Покупка домой", "Фермерский рынок"}, weight: 28},
		{categoryKey: "food", subcategoryKey: "cafe", accountKey: "card", minAmount: 450, maxAmount: 3200, comments: []string{"Кофе", "Обед", "Ужин в кафе", "Доставка еды"}, weight: 18},
		{categoryKey: "food", subcategoryKey: "fastfood", accountKey: "cash", minAmount: 250, maxAmount: 1400, comments: []string{"Перекус", "Фастфуд", "Снэки"}, weight: 11},
		{categoryKey: "transport", subcategoryKey: "taxi", accountKey: "card", minAmount: 220, maxAmount: 1800, comments: []string{"Поездка на такси", "Такси до работы", "Такси ночью"}, weight: 16},
		{categoryKey: "transport", subcategoryKey: "public", accountKey: "card", minAmount: 55, maxAmount: 290, comments: []string{"Метро", "Автобус", "Проезд"}, weight: 10},
		{categoryKey: "home", accountKey: "card", minAmount: 1200, maxAmount: 9200, comments: []string{"Коммунальные услуги", "Интернет", "Оплата дома"}, weight: 8},
		{categoryKey: "fun", accountKey: "card", minAmount: 400, maxAmount: 6500, comments: []string{"Кино", "Подписки", "Развлечения"}, weight: 9},
		{categoryKey: "health", subcategoryKey: "pharmacy", accountKey: "card", minAmount: 250, maxAmount: 4800, comments: []string{"Аптека", "Лекарства", "Витамины"}, weight: 6},
		{categoryKey: "shopping", subcategoryKey: "marketplace", accountKey: "card", minAmount: 500, maxAmount: 9000, comments: []string{"Покупка на маркетплейсе", "Одежда", "Товары для дома"}, weight: 12},
		{categoryKey: "shopping", accountKey: "card", minAmount: 300, maxAmount: 4500, comments: []string{"Покупки", "Магазин", "Хозтовары"}, weight: 7},
	}

	incomePatterns := []txPattern{
		{categoryKey: "income", accountKey: "card", minAmount: 45000, maxAmount: 120000, comments: []string{"Зарплата"}, weight: 8},
		{categoryKey: "income", accountKey: "savings", minAmount: 8000, maxAmount: 45000, comments: []string{"Подработка", "Фриланс", "Бонус"}, weight: 5},
		{categoryKey: "income", accountKey: "cash", minAmount: 1500, maxAmount: 12000, comments: []string{"Возврат долга", "Кэшбэк", "Продажа вещей"}, weight: 3},
	}

	weightedPick := func(patterns []txPattern) txPattern {
		total := 0
		for _, p := range patterns {
			total += p.weight
		}
		roll := rng.Intn(total)
		for _, p := range patterns {
			roll -= p.weight
			if roll < 0 {
				return p
			}
		}
		return patterns[len(patterns)-1]
	}

	randAmount := func(minVal, maxVal float64) float64 {
		if maxVal <= minVal {
			return minVal
		}
		value := minVal + rng.Float64()*(maxVal-minVal)
		return float64(int(value*100+0.5)) / 100
	}

	randComment := func(pattern txPattern) string {
		if len(pattern.comments) == 0 {
			return ""
		}
		return pattern.comments[rng.Intn(len(pattern.comments))]
	}

	buildDateTime := func(day time.Time) time.Time {
		hour := 7 + rng.Intn(16)
		minute := rng.Intn(60)
		return time.Date(day.Year(), day.Month(), day.Day(), hour, minute, 0, 0, day.Location())
	}

	start := time.Now().AddDate(0, 0, -90)
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	today := time.Now()

	hotDays := map[string]int{
		formatDay(today.AddDate(0, 0, -1)): 26,
		formatDay(today.AddDate(0, 0, -2)): 21,
		formatDay(today.AddDate(0, 0, -6)): 18,
		formatDay(today.AddDate(0, 0, -12)): 23,
	}

	createdTransactions := 0
	for day := start; !day.After(today); day = day.AddDate(0, 0, 1) {
		expenseCount := 2 + rng.Intn(4)
		if day.Weekday() == time.Saturday || day.Weekday() == time.Sunday {
			expenseCount += 2 + rng.Intn(3)
		}
		if extra, ok := hotDays[formatDay(day)]; ok {
			expenseCount += extra
		}

		for i := 0; i < expenseCount; i++ {
			pattern := weightedPick(expensePatterns)
			acc := accountsByKey[pattern.accountKey]
			cat := categoriesByKey[pattern.categoryKey]
			amount := randAmount(pattern.minAmount, pattern.maxAmount)

			tx := Transaction{
				AccountID:  acc.ID,
				CategoryID: cat.ID,
				Type:       Expense,
				Amount:     amount,
				Comment:    randComment(pattern),
				Date:       buildDateTime(day),
			}
			if pattern.subcategoryKey != "" {
				sub := subcategoriesByKey[pattern.subcategoryKey]
				tx.SubCategoryID = sub.ID
			}

			if err := AddTransaction(db, tx); err != nil {
				return fmt.Errorf("failed to add expense transaction: %w", err)
			}
			createdTransactions++
		}

		incomeCount := 0
		if day.Day() == 5 || day.Day() == 20 {
			incomeCount = 1 + rng.Intn(2)
		} else if rng.Intn(100) < 10 {
			incomeCount = 1
		}

		for i := 0; i < incomeCount; i++ {
			pattern := weightedPick(incomePatterns)
			acc := accountsByKey[pattern.accountKey]
			cat := categoriesByKey[pattern.categoryKey]
			amount := randAmount(pattern.minAmount, pattern.maxAmount)

			tx := Transaction{
				AccountID:  acc.ID,
				CategoryID: cat.ID,
				Type:       Income,
				Amount:     amount,
				Comment:    randComment(pattern),
				Date:       buildDateTime(day),
			}

			if err := AddTransaction(db, tx); err != nil {
				return fmt.Errorf("failed to add income transaction: %w", err)
			}
			createdTransactions++
		}
	}

	for i := 0; i < 8; i++ {
		day := today.AddDate(0, 0, -rng.Intn(60))

		var card Account
		if err := db.First(&card, accountsByKey["card"].ID).Error; err != nil {
			return fmt.Errorf("failed to load card account for transfer seed: %w", err)
		}

		var savings Account
		if err := db.First(&savings, accountsByKey["savings"].ID).Error; err != nil {
			return fmt.Errorf("failed to load savings account for transfer seed: %w", err)
		}

		fromID := card.ID
		toID := savings.ID
		comment := "Перевод на сбережения"
		available := card.Balance

		if available < 1500 {
			fromID = savings.ID
			toID = card.ID
			comment = "Перевод с накоплений"
			available = savings.Balance
		}

		if available < 1500 {
			continue
		}

		maxTransfer := available * 0.35
		if maxTransfer > 25000 {
			maxTransfer = 25000
		}
		if maxTransfer < 1500 {
			continue
		}

		transfer := TransferTransaction{
			AccountID:         fromID,
			TransferAccountID: toID,
			Type:              Transfer,
			Amount:            randAmount(1500, maxTransfer),
			Comment:           comment,
			Date:              buildDateTime(day),
		}

		if err := TransferToAnotherAcc(db, transfer); err != nil {
			return fmt.Errorf("failed to add transfer transaction: %w", err)
		}
	}

	log.Printf("GenerateTestData: generated realistic demo data, transactions=%d", createdTransactions)
	return nil
}

func ensureUniqueName(db *gorm.DB, model interface{}, base string) (string, error) {
	for i := 0; i < 1000; i++ {
		candidate := base
		if i > 0 {
			candidate = fmt.Sprintf("%s %d", base, i+1)
		}

		var count int64
		if err := db.Model(model).Where("name = ?", candidate).Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("could not generate unique name for %s", base)
}

func formatDay(day time.Time) string {
	return day.Format("2006-01-02")
}
