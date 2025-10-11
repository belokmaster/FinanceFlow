package database

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

func AddCategory(db *gorm.DB, cat Category) error {
	log.Printf("Adding new category: %s", cat.Name)

	if cat.Name == "" {
		log.Printf("Error: empty category name")
		return fmt.Errorf("empty name of category")
	}

	result := db.Create(&cat)
	if result.Error != nil {
		log.Printf("Error creating category %s: %v", cat.Name, result.Error)
		return fmt.Errorf("problem to create a new category: %v", result.Error)
	}

	log.Printf("Category created successfully: ID=%d, Name=%s", cat.ID, cat.Name)
	return nil
}

func DeleteCategory(db *gorm.DB, id int) error {
	log.Printf("Deleting category with ID: %d", id)

	result := db.Delete(&Category{}, id)
	if result.Error != nil {
		log.Printf("Error deleting category ID %d: %v", id, result.Error)
		return fmt.Errorf("problem with delete category in db: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		log.Printf("Category with ID %d not found for deletion", id)
		return fmt.Errorf("category with id %d not found", id)
	}

	log.Printf("Category deleted successfully: ID=%d", id)
	return nil
}
