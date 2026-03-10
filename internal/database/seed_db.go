package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

func GenerateTestData(db *gorm.DB) error {
	suffix := time.Now().Format("20060102150405")

	cashAcc := Account{
		Name:         fmt.Sprintf("Demo Cash %s", suffix),
		Balance:      50000,
		CurrencyCode: Ruble,
		Color:        "4cd67a",
		IconCode:     Coin,
	}
	cardAcc := Account{
		Name:         fmt.Sprintf("Demo Card %s", suffix),
		Balance:      120000,
		CurrencyCode: Ruble,
		Color:        "4aa8ff",
		IconCode:     Card,
	}
	savingsAcc := Account{
		Name:         fmt.Sprintf("Demo Savings %s", suffix),
		Balance:      300000,
		CurrencyCode: Ruble,
		Color:        "ffb84d",
		IconCode:     Mark,
	}

	if err := CreateNewAccount(db, cashAcc); err != nil {
		return fmt.Errorf("failed to create demo cash account: %w", err)
	}
	if err := CreateNewAccount(db, cardAcc); err != nil {
		return fmt.Errorf("failed to create demo card account: %w", err)
	}
	if err := CreateNewAccount(db, savingsAcc); err != nil {
		return fmt.Errorf("failed to create demo savings account: %w", err)
	}

	var accounts []Account
	if err := db.Where("name IN ?", []string{cashAcc.Name, cardAcc.Name, savingsAcc.Name}).Find(&accounts).Error; err != nil {
		return fmt.Errorf("failed to get created accounts: %w", err)
	}
	if len(accounts) != 3 {
		return fmt.Errorf("expected 3 created accounts, got %d", len(accounts))
	}

	accByName := make(map[string]Account, 3)
	for _, acc := range accounts {
		accByName[acc.Name] = acc
	}

	salaryCat := Category{
		Name:     fmt.Sprintf("Demo Salary %s", suffix),
		Color:    "4cd67a",
		IconCode: Coin_cat,
	}
	foodCat := Category{
		Name:     fmt.Sprintf("Demo Food %s", suffix),
		Color:    "ff7f50",
		IconCode: Food,
	}
	transportCat := Category{
		Name:     fmt.Sprintf("Demo Transport %s", suffix),
		Color:    "5f9ea0",
		IconCode: Bus1,
	}

	if err := CreateCategory(db, salaryCat); err != nil {
		return fmt.Errorf("failed to create salary category: %w", err)
	}
	if err := CreateCategory(db, foodCat); err != nil {
		return fmt.Errorf("failed to create food category: %w", err)
	}
	if err := CreateCategory(db, transportCat); err != nil {
		return fmt.Errorf("failed to create transport category: %w", err)
	}

	var categories []Category
	if err := db.Where("name IN ?", []string{salaryCat.Name, foodCat.Name, transportCat.Name}).Find(&categories).Error; err != nil {
		return fmt.Errorf("failed to get created categories: %w", err)
	}
	if len(categories) != 3 {
		return fmt.Errorf("expected 3 created categories, got %d", len(categories))
	}

	catByName := make(map[string]Category, 3)
	for _, cat := range categories {
		catByName[cat.Name] = cat
	}

	groceriesSub := SubCategory{
		Name:       fmt.Sprintf("Demo Groceries %s", suffix),
		Color:      "ff9f43",
		IconCode:   Cooking1,
		CategoryID: catByName[foodCat.Name].ID,
	}
	cafeSub := SubCategory{
		Name:       fmt.Sprintf("Demo Cafe %s", suffix),
		Color:      "ee5253",
		IconCode:   Restaurant1,
		CategoryID: catByName[foodCat.Name].ID,
	}

	if err := AddSubCategory(db, groceriesSub); err != nil {
		return fmt.Errorf("failed to create groceries subcategory: %w", err)
	}
	if err := AddSubCategory(db, cafeSub); err != nil {
		return fmt.Errorf("failed to create cafe subcategory: %w", err)
	}

	var subCategories []SubCategory
	if err := db.Where("name IN ?", []string{groceriesSub.Name, cafeSub.Name}).Find(&subCategories).Error; err != nil {
		return fmt.Errorf("failed to get created subcategories: %w", err)
	}
	if len(subCategories) != 2 {
		return fmt.Errorf("expected 2 created subcategories, got %d", len(subCategories))
	}

	subByName := make(map[string]SubCategory, 2)
	for _, sub := range subCategories {
		subByName[sub.Name] = sub
	}

	now := time.Now()
	demoTransactions := []Transaction{
		{
			AccountID:  accByName[cardAcc.Name].ID,
			CategoryID: catByName[salaryCat.Name].ID,
			Type:       Income,
			Amount:     90000,
			Comment:    "Demo salary",
			Date:       now.AddDate(0, 0, -12),
		},
		{
			AccountID:     accByName[cashAcc.Name].ID,
			CategoryID:    catByName[foodCat.Name].ID,
			SubCategoryID: subByName[groceriesSub.Name].ID,
			Type:          Expense,
			Amount:        2300,
			Comment:       "Demo groceries",
			Date:          now.AddDate(0, 0, -10),
		},
		{
			AccountID:     accByName[cardAcc.Name].ID,
			CategoryID:    catByName[foodCat.Name].ID,
			SubCategoryID: subByName[cafeSub.Name].ID,
			Type:          Expense,
			Amount:        1800,
			Comment:       "Demo cafe",
			Date:          now.AddDate(0, 0, -9),
		},
		{
			AccountID:  accByName[cardAcc.Name].ID,
			CategoryID: catByName[transportCat.Name].ID,
			Type:       Expense,
			Amount:     950,
			Comment:    "Demo taxi",
			Date:       now.AddDate(0, 0, -7),
		},
		{
			AccountID:     accByName[cashAcc.Name].ID,
			CategoryID:    catByName[foodCat.Name].ID,
			SubCategoryID: subByName[groceriesSub.Name].ID,
			Type:          Expense,
			Amount:        1200,
			Comment:       "Demo market",
			Date:          now.AddDate(0, 0, -5),
		},
		{
			AccountID:  accByName[savingsAcc.Name].ID,
			CategoryID: catByName[salaryCat.Name].ID,
			Type:       Income,
			Amount:     30000,
			Comment:    "Demo freelance",
			Date:       now.AddDate(0, 0, -3),
		},
	}

	for i, tr := range demoTransactions {
		if err := AddTransaction(db, tr); err != nil {
			return fmt.Errorf("failed to add demo transaction #%d: %w", i+1, err)
		}
	}

	demoTransfer := TransferTransaction{
		AccountID:         accByName[cardAcc.Name].ID,
		TransferAccountID: accByName[savingsAcc.Name].ID,
		Type:              Transfer,
		Amount:            10000,
		Comment:           "Demo transfer to savings",
		Date:              now.AddDate(0, 0, -2),
	}
	if err := TransferToAnotherAcc(db, demoTransfer); err != nil {
		return fmt.Errorf("failed to add demo transfer: %w", err)
	}

	log.Printf("GenerateTestData: demo data generated successfully with suffix %s", suffix)
	return nil
}
