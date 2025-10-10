package database

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

func InitDB(db *gorm.DB) error {
	log.Println("Starting database migration...")

	err := db.AutoMigrate(&Account{}, &Category{}, &SubCategory{}, &Transaction{}, &TransferTransaction{})
	if err != nil {
		log.Printf("Error creating database: %v", err)
		log.Fatalf("problem to create gorm bd: %v", err)
	}

	log.Println("The database created successfully")
	return nil
}

func GetAccounts(db *gorm.DB) ([]Account, error) {
	log.Println("Getting all accounts from database")

	var accounts []Account
	result := db.Find(&accounts)

	if result.Error != nil {
		log.Printf("Error getting accounts: %v", result.Error)
		return nil, fmt.Errorf("problem with get accounts in db: %v", result.Error)
	}

	log.Printf("Successfully got %d accounts", len(accounts))
	return accounts, nil
}
