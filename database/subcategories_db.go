package database

import (
	"fmt"

	"gorm.io/gorm"
)

func AddSubCategory(db *gorm.DB, subCat SubCategory) error {
	if subCat.Name == "" {
		return fmt.Errorf("empty name of category")
	}

	var cat Category
	err := db.First(&cat, subCat.CategoryID).Error
	if err != nil {
		return fmt.Errorf("category with ID %d not found", subCat.CategoryID)
	}

	result := db.Create(&subCat)
	if result.Error != nil {
		return fmt.Errorf("problem to create a new sub_category: %v", result.Error)
	}

	return nil
}

func DeleteSubCategory(db *gorm.DB, id int) error {
	result := db.Delete(&SubCategory{}, id)
	if result.Error != nil {
		return fmt.Errorf("problem with delete sub_category in db: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("sub_category with id %d not found", id)
	}

	return nil
}
