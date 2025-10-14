package handlers

import (
	"finance_flow/internal/database"
	"finance_flow/internal/icons"
	"html/template"
	"log"

	"gorm.io/gorm"
)

func getAccountsForView(db *gorm.DB) ([]AccountView, error) {
	log.Printf("getAccountsForView: Getting accounts from database")
	accounts, err := database.GetAccounts(db)
	if err != nil {
		log.Printf("getAccountsForView: Error getting accounts: %v", err)
		return nil, err
	}

	log.Printf("getAccountsForView: Retrieved %d accounts from database", len(accounts))

	var accountsForView []AccountView
	for _, acc := range accounts {
		symbol, ok := database.CurrencySymbols[acc.CurrencyCode]
		if !ok {
			log.Printf("getAccountsForView: Unknown currency code %d for account %s, using default symbol", acc.CurrencyCode, acc.Name)
			symbol = "?"
		}

		color := "#" + acc.Color
		if color == "" {
			log.Printf("getAccountsForView: Empty color for account %s, using default", acc.Name)
			color = "#4cd67a"
		}

		iconFileName, ok := database.IconAccountFiles[acc.IconCode]
		if !ok {
			log.Printf("getAccountsForView: Unknown icon code %d for account %s, using default icon", acc.IconCode, acc.Name)
			iconFileName = "Wallet"
		}

		iconHTML, ok := icons.AccountIconCache[iconFileName]
		if !ok {
			log.Printf("getAccountsForView: Icon file %s not found in cache for account %s, using coin icon", iconFileName, acc.Name)
			iconHTML = icons.AccountIconCache["coin"]
		}

		accountsForView = append(accountsForView, AccountView{
			ID:             acc.ID,
			Name:           acc.Name,
			Balance:        acc.Balance,
			CurrencySymbol: symbol,
			Color:          color,
			IconKey:        iconFileName,
			IconHTML:       iconHTML,
		})
	}

	return accountsForView, nil
}

func convertSubcategoriesToView(subCategories []database.SubCategory) map[uint][]SubCategoryView {
	subcategoriesByParent := make(map[uint][]SubCategoryView)

	for _, subCat := range subCategories {
		color := "#" + subCat.Color
		if color == "" {
			color = "#4cd67a"
		}

		iconFileName, ok := database.IconSubCategoryFiles[subCat.IconCode]
		if !ok {
			iconFileName = "Restaraunt1"
		}

		iconHTML, ok := icons.SubCategoryIconCache[iconFileName]
		if !ok {
			iconHTML = icons.SubCategoryIconCache["Restaraunt1"]
		}

		subcategoryView := SubCategoryView{
			ID:       subCat.ID,
			Name:     subCat.Name,
			Color:    color,
			IconKey:  iconFileName,
			IconHTML: iconHTML,
			ParentID: subCat.CategoryID,
		}

		subcategoriesByParent[subCat.CategoryID] = append(subcategoriesByParent[subCat.CategoryID], subcategoryView)
	}

	return subcategoriesByParent
}

func convertCategoriesToView(categories []database.Category, subcategoriesByParent map[uint][]SubCategoryView) []CategoryView {
	var categoriesForView []CategoryView

	for _, cat := range categories {
		color := "#" + cat.Color
		if color == "" {
			color = "#4cd67a"
		}

		iconFileName, ok := database.IconCategoryFiles[cat.IconCode]
		if !ok {
			iconFileName = "Food"
		}

		iconHTML, ok := icons.CategoryIconCache[iconFileName]
		if !ok {
			iconHTML = icons.CategoryIconCache["Food"]
		}

		subcats := subcategoriesByParent[cat.ID]
		if subcats == nil {
			subcats = []SubCategoryView{}
		}

		categoriesForView = append(categoriesForView, CategoryView{
			ID:            cat.ID,
			Name:          cat.Name,
			Color:         color,
			IconKey:       iconFileName,
			IconHTML:      iconHTML,
			Subcategories: subcats,
		})
	}

	return categoriesForView
}

func getTransactionsForView(db *gorm.DB) ([]TransactionView, error) {
	log.Printf("getTransactionsForView: Getting transactions from database")

	transactions, err := database.GetTransactions(db)
	if err != nil {
		log.Printf("getTransactionsForView: Error getting transactions: %v", err)
		return nil, err
	}

	log.Printf("getTransactionsForView: Retrieved %d transactions from database", len(transactions))

	var transactionsForView []TransactionView
	for _, tx := range transactions {
		symbol, ok := database.CurrencySymbols[tx.Account.CurrencyCode]
		if !ok {
			log.Printf("getTransactionsForView: Unknown currency code %d for transaction %d, using default symbol", tx.Account.CurrencyCode, tx.ID)
			symbol = "?"
		}

		formattedDate := tx.Date.Format("02.01.2006")
		dateOnly := tx.Date.Format("2006-01-02")

		var (
			displayColor       string
			displayIconHTML    template.HTML
			displayName        string
			parentCategoryName string
		)

		if tx.SubCategoryID != 0 && tx.SubCategory.ID != 0 {
			log.Printf("Transaction %d has subcategory: %s", tx.ID, tx.SubCategory.Name)

			displayColor = tx.SubCategory.Color
			if displayColor == "" {
				displayColor = "4cd67a"
			}

			subCategoryIconFileName, ok := database.IconSubCategoryFiles[tx.SubCategory.IconCode]
			if !ok {
				subCategoryIconFileName = "Restaraunt1"
			}

			displayIconHTML, ok = icons.SubCategoryIconCache[subCategoryIconFileName]
			if !ok {
				displayIconHTML = icons.SubCategoryIconCache["Restaraunt1"]
			}

			displayName = tx.SubCategory.Name
			parentCategoryName = tx.Category.Name

		} else {
			displayColor = tx.Category.Color
			if displayColor == "" {
				displayColor = "4cd67a"
			}

			categoryIconFileName, ok := database.IconCategoryFiles[tx.Category.IconCode]
			if !ok {
				categoryIconFileName = "Food"
			}

			displayIconHTML, ok = icons.CategoryIconCache[categoryIconFileName]
			if !ok {
				displayIconHTML = icons.CategoryIconCache["Food"]
			}

			displayName = tx.Category.Name
		}

		transactionView := TransactionView{
			ID:                 tx.ID,
			Type:               tx.Type,
			Amount:             tx.Amount,
			AccountID:          tx.AccountID,
			AccountName:        tx.Account.Name,
			CurrencySymbol:     symbol,
			CategoryID:         tx.CategoryID,
			CategoryName:       tx.Category.Name,
			CategoryColor:      displayColor,
			CategoryIconHTML:   displayIconHTML,
			DisplayName:        displayName,
			ParentCategoryName: parentCategoryName,
			Date:               dateOnly,
			FormattedDate:      formattedDate,
			Description:        tx.Comment,
		}

		if tx.SubCategoryID != 0 && tx.SubCategory.ID != 0 {
			transactionView.SubCategoryID = &tx.SubCategoryID
			subCatName := tx.SubCategory.Name
			transactionView.SubCategoryName = &subCatName
		}

		transactionsForView = append(transactionsForView, transactionView)
	}

	return transactionsForView, nil
}
