package database

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

func InitDB(db *gorm.DB) error {
	err := db.AutoMigrate(&Account{}, &Category{}, &SubCategory{}, &Transaction{}, &TransferTransaction{})
	if err != nil {
		log.Fatalf("problem to create gorm bd: %v", err)
	}

	log.Println("The database created successfully")
	return nil
}

func GetAccounts(db *gorm.DB) ([]Account, error) {
	var accounts []Account
	result := db.Find(&accounts)

	if result.Error != nil {
		return nil, fmt.Errorf("problem with get accounts in db: %v", result.Error)
	}

	return accounts, nil
}
