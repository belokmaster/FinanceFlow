package database

import (
	"fmt"

	"gorm.io/gorm"
)

func AddCategory(db *gorm.DB, cat Category) error {
	if cat.Name == "" {
		return fmt.Errorf("empty name of category")
	}

	result := db.Create(&cat)
	if result.Error != nil {
		return fmt.Errorf("problem to create a new category: %v", result.Error)
	}

	return nil
}

func DeleteCategory(db *gorm.DB, id int) error {
	result := db.Delete(&Category{}, id)
	if result.Error != nil {
		return fmt.Errorf("problem with delete category in db: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("category with id %d not found", id)
	}

	return nil
}
