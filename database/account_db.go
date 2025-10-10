package database

import (
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
)

func CreateNewAccount(db *gorm.DB, acc Account) error {
	log.Printf("Creating new account: %s", acc.Name)

	if acc.Name == "" {
		log.Printf("Error: empty account name")
		return fmt.Errorf("empty name of account")
	}

	result := db.Create(&acc)
	if result.Error != nil {
		log.Printf("Error creating account %s: %v", acc.Name, result.Error)
		return fmt.Errorf("problem to create a new acc: %v", result.Error)
	}

	return nil
}

func DeleteAccount(db *gorm.DB, id int) error {
	log.Printf("Deleting account with ID: %d", id)

	result := db.Delete(&Account{}, id)
	if result.Error != nil {
		log.Printf("Error deleting account ID %d: %v", id, result.Error)
		return fmt.Errorf("problem with delete account in db: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		log.Printf("Account with ID %d not found for deletion", id)
		return fmt.Errorf("account with id %d not found", id)
	}

	log.Printf("Account deleted successfully ID: %d", id)
	return nil
}

func ChangeAccountColor(db *gorm.DB, id int, newColor string) error {
	log.Printf("Changing color for account ID %d to: %s", id, newColor)

	if newColor == "" {
		log.Printf("Error: empty color for account ID %d", id)
		return fmt.Errorf("empty color")
	}

	newColor = strings.TrimPrefix(newColor, "#")

	var acc Account
	err := db.First(&acc, id).Error
	if err != nil {
		log.Printf("Error finding account ID %d for color change: %v", id, err)
		return fmt.Errorf("account with ID %d not found", id)
	}

	err = db.Model(&Account{}).Where("id = ?", id).Update("Color", newColor).Error
	if err != nil {
		log.Printf("Error updating color for account ID %d: %v", id, err)
		return fmt.Errorf("problem with change color in db: %v", err)
	}

	log.Printf("Color updated successfully for account ID %d", id)
	return nil
}

func ChangeAccountIcon(db *gorm.DB, id int, icon_id int) error {
	log.Printf("Changing icon for account ID %d to icon ID: %d", id, icon_id)

	if _, ok := IconFiles[TypeIcons(icon_id)]; !ok {
		log.Printf("Error: icon ID %d does not exist", icon_id)
		return fmt.Errorf("icon with ID %d does not exist", icon_id)
	}

	var acc Account
	err := db.First(&acc, id).Error
	if err != nil {
		log.Printf("Error finding account ID %d for icon change: %v", id, err)
		return fmt.Errorf("account with ID %d not found", id)
	}

	log.Printf("Found account: ID=%d, Name=%s, Old icon=%d, New icon=%d", acc.ID, acc.Name, acc.IconCode, icon_id)

	result := db.Model(&Account{}).Where("id = ?", id).Update("icon_code", TypeIcons(icon_id))
	if result.Error != nil {
		log.Printf("Error updating icon for account ID %d: %v", id, result.Error)
		return fmt.Errorf("problem with changing icon in db: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		log.Printf("No rows affected when updating icon for account ID %d", id)
		return fmt.Errorf("no rows were updated - account may not exist")
	}

	log.Printf("Icon updated successfully for account ID %d", id)
	return nil
}

func ChangeAccountName(db *gorm.DB, id int, newName string) error {
	var acc Account
	err := db.First(&acc, id).Error
	if err != nil {
		return fmt.Errorf("account with ID %d not found", id)
	}

	result := db.Model(&Account{}).Where("id = ?", id).Update("Name", newName)
	if result.Error != nil {
		return fmt.Errorf("problem with change name in db: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows were updated - account may not exist")
	}

	return nil
}
