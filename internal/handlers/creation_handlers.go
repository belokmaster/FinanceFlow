package handlers

import (
	"finance_flow/internal/database"
	"finance_flow/internal/icons"
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
