package database

import (
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
)

func CreateCategory(db *gorm.DB, cat Category) error {
	if strings.TrimSpace(cat.Name) == "" {
		err := fmt.Errorf("category name cannot be empty")
		log.Println("Error creating category:", err)
		return err
	}

	log.Printf("Adding new category: %s", cat.Name)

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

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("Error starting transaction: %v", tx.Error)
		return fmt.Errorf("problem starting transaction: %v", tx.Error)
	}

	var subCategories []SubCategory
	if err := tx.Where("category_id = ?", id).Find(&subCategories).Error; err != nil {
		tx.Rollback()
		log.Printf("Error finding subcategories for category ID %d: %v", id, err)
		return fmt.Errorf("problem finding subcategories: %v", err)
	}

	for _, subCategory := range subCategories {
		if err := DeleteSubCategory(tx, int(subCategory.ID)); err != nil {
			tx.Rollback()
			log.Printf("Error deleting subcategory ID %d for category ID %d: %v", subCategory.ID, id, err)
			return fmt.Errorf("problem deleting subcategory %d: %v", subCategory.ID, err)
		}
	}
	log.Printf("Deleted %d subcategories for category ID %d", len(subCategories), id)

	var transactions []Transaction
	if err := tx.Where("category_id = ?", id).Find(&transactions).Error; err != nil {
		tx.Rollback()
		log.Printf("Error finding direct transactions for category ID %d: %v", id, err)
		return fmt.Errorf("problem finding direct transactions: %v", err)
	}

	for _, transaction := range transactions {
		if err := DeleteTransaction(tx, int(transaction.ID)); err != nil {
			tx.Rollback()
			log.Printf("Error deleting direct transaction ID %d for category ID %d: %v", transaction.ID, id, err)
			return fmt.Errorf("problem deleting direct transaction %d: %v", transaction.ID, err)
		}
	}
	log.Printf("Deleted %d direct transactions for category ID %d", len(transactions), id)

	result := tx.Delete(&Category{}, id)
	if result.Error != nil {
		tx.Rollback()
		log.Printf("Error deleting category ID %d: %v", id, result.Error)
		return fmt.Errorf("problem with delete category in db: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		log.Printf("Category with ID %d not found for deletion", id)
		return fmt.Errorf("category with id %d not found", id)
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction for category ID %d: %v", id, err)
		return fmt.Errorf("problem committing transaction: %v", err)
	}

	log.Printf("Category ID %d and all related subcategories/transactions deleted successfully", id)
	return nil
}

func ChangeCategoryColor(db *gorm.DB, id int, newColor string) error {
	log.Printf("Changing color for category ID %d to: %s", id, newColor)

	if newColor == "" {
		log.Printf("Error: empty color for category ID %d", id)
		return fmt.Errorf("empty color")
	}

	newColor = strings.TrimPrefix(newColor, "#")

	result := db.Model(&Category{}).Where("id = ?", id).Update("Color", newColor)
	if result.Error != nil {
		log.Printf("Error updating color for category ID %d: %v", id, result.Error)
		return fmt.Errorf("problem with change color in db: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		log.Printf("No rows affected when updating color for category ID %d", id)
		return fmt.Errorf("no rows were updated - category may not exist")
	}

	log.Printf("Color updated successfully for category ID %d", id)
	return nil
}

func ChangeCategoryIcon(db *gorm.DB, id int, icon_id int) error {
	log.Printf("Changing icon for category ID %d to icon ID: %d", id, icon_id)

	if _, ok := IconCategoryFiles[TypeCategoryIcons(icon_id)]; !ok {
		log.Printf("Error: icon ID %d does not exist", icon_id)
		return fmt.Errorf("icon with ID %d does not exist", icon_id)
	}

	result := db.Model(&Category{}).Where("id = ?", id).Update("icon_code", TypeCategoryIcons(icon_id))
	if result.Error != nil {
		log.Printf("Error updating icon for category ID %d: %v", id, result.Error)
		return fmt.Errorf("problem with changing icon in db: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		log.Printf("No rows affected when updating category for category ID %d", id)
		return fmt.Errorf("no rows were updated - category may not exist")
	}

	log.Printf("Icon updated successfully for category ID %d", id)
	return nil
}

func ChangeCategoryName(db *gorm.DB, id int, newName string) error {
	newName = strings.TrimSpace(newName)
	if newName == "" {
		return fmt.Errorf("new name cannot be empty")
	}

	result := db.Model(&Category{}).Where("id = ?", id).Update("Name", newName)
	if result.Error != nil {
		return fmt.Errorf("problem with change name in db: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows were updated - account may not exist")
	}

	log.Printf("Name updated successfully for category ID %d", id)
	return nil
}
