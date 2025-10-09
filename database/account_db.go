package database

import (
	"fmt"

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
