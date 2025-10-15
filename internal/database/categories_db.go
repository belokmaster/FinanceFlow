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
