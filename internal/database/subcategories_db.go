package database

import (
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
)

func AddSubCategory(db *gorm.DB, subCat SubCategory) error {
	if strings.TrimSpace(subCat.Name) == "" {
		err := fmt.Errorf("subcategory name cannot be empty")
		log.Println("Error creating subcategory:", err)
		return err
	}

	result := db.Create(&subCat)
	if result.Error != nil {
		log.Printf("Error creating subcategory %s: %v", subCat.Name, result.Error)
		return fmt.Errorf("problem to create a new subcategory: %v", result.Error)
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

func ChangeSubCategoryColor(db *gorm.DB, id int, newColor string) error {
	log.Printf("Changing color for sub_category ID %d to: %s", id, newColor)

	if newColor == "" {
		log.Printf("Error: empty color for sub_category ID %d", id)
		return fmt.Errorf("empty color")
	}

	newColor = strings.TrimPrefix(newColor, "#")

	result := db.Model(&SubCategory{}).Where("id = ?", id).Update("Color", newColor)
	if result.Error != nil {
		log.Printf("Error updating color for sub_category ID %d: %v", id, result.Error)
		return fmt.Errorf("problem with change color in db: %v", result.Error)
	}

	log.Printf("Color updated successfully for sub_category ID %d", id)
	return nil
}

func ChangeSubCategoryIcon(db *gorm.DB, id int, icon_id int) error {
	log.Printf("Changing icon for sub_category ID %d to icon ID: %d", id, icon_id)

	if _, ok := IconSubCategoryFiles[TypeSubCategoryIcons(icon_id)]; !ok {
		log.Printf("Error: icon ID %d does not exist", icon_id)
		return fmt.Errorf("icon with ID %d does not exist", icon_id)
	}

	result := db.Model(&SubCategory{}).Where("id = ?", id).Update("icon_code", TypeSubCategoryIcons(icon_id))
	if result.Error != nil {
		log.Printf("Error updating icon for sub_category ID %d: %v", id, result.Error)
		return fmt.Errorf("problem with changing icon in db: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		log.Printf("No rows affected when updating sub_category for sub_category ID %d", id)
		return fmt.Errorf("no rows were updated - sub_category may not exist")
	}

	log.Printf("Icon updated successfully for sub_category ID %d", id)
	return nil
}

func ChangeSubCategoryName(db *gorm.DB, id int, newName string) error {
	newName = strings.TrimSpace(newName)
	if newName == "" {
		return fmt.Errorf("new name cannot be empty")
	}

	result := db.Model(&SubCategory{}).Where("id = ?", id).Update("Name", newName)
	if result.Error != nil {
		return fmt.Errorf("problem with change name in db: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows were updated - account may not exist")
	}

	return nil
}
