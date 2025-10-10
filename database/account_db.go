package database

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func CreateNewAccount(db *gorm.DB, acc Account) error {
	if acc.Name == "" {
		return fmt.Errorf("empty name of account")
	}

	result := db.Create(&acc)
	if result.Error != nil {
		return fmt.Errorf("problem to create a new acc: %v", result.Error)
	}

	return nil
}

func DeleteAccount(db *gorm.DB, id int) error {
	result := db.Delete(&Account{}, id)
	if result.Error != nil {
		return fmt.Errorf("problem with delete account in db: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("account with id %d not found", id)
	}

	return nil
}

func ChangeAccountColor(db *gorm.DB, id int, newColor string) error {
	if newColor == "" {
		return fmt.Errorf("empty color")
	}

	newColor = strings.TrimPrefix(newColor, "#")

	var acc Account
	err := db.First(&acc, id).Error
	if err != nil {
		return fmt.Errorf("account with ID %d not found", id)
	}

	err = db.Model(&Account{}).Where("id = ?", id).Update("Color", newColor).Error
	if err != nil {
		return fmt.Errorf("problem with change color in db: %v", err)
	}

	return nil
}
