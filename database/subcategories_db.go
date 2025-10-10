package database

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

func AddSubCategory(db *gorm.DB, subCat SubCategory) error {
	log.Printf("Adding new subcategory: %s for category ID: %d", subCat.Name, subCat.CategoryID)

	if subCat.Name == "" {
		log.Printf("Error: empty subcategory name")
		return fmt.Errorf("empty name of category")
	}

	var cat Category
	err := db.First(&cat, subCat.CategoryID).Error
	if err != nil {
		log.Printf("Error: category with ID %d not found", subCat.CategoryID)
		return fmt.Errorf("category with ID %d not found", subCat.CategoryID)
	}

	log.Printf("Found parent category: ID=%d, Name=%s", cat.ID, cat.Name)

	result := db.Create(&subCat)
	if result.Error != nil {
		log.Printf("Error creating subcategory %s: %v", subCat.Name, result.Error)
		return fmt.Errorf("problem to create a new sub_category: %v", result.Error)
	}

	log.Printf("Subcategory created successfully: ID=%d, Name=%s, Parent Category ID=%d",
		subCat.ID, subCat.Name, subCat.CategoryID)
	return nil
}

func DeleteSubCategory(db *gorm.DB, id int) error {
	log.Printf("Deleting subcategory with ID: %d", id)

	result := db.Delete(&SubCategory{}, id)
	if result.Error != nil {
		log.Printf("Error deleting subcategory ID %d: %v", id, result.Error)
		return fmt.Errorf("problem with delete sub_category in db: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		log.Printf("Subcategory with ID %d not found for deletion", id)
		return fmt.Errorf("sub_category with id %d not found", id)
	}

	log.Printf("Subcategory deleted successfully: ID=%d", id)
	return nil
}
