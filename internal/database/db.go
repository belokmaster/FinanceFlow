package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDatabase() *gorm.DB {
	dbPath := "internal/database/database.db" // mb change this point later

	// check if db exists
	_, err := os.Stat(dbPath)
	dbExists := !os.IsNotExist(err)

	// open the db using gorm
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("problem connecting to db: %v", err)
	}

	// initialize a new db if it doesnt exist
	if !dbExists {
		log.Println("initialize a new db")
		err := migrateDB(db)
		if err != nil {
			log.Fatalf("failed to initialize the db: %v", err)
		}
	}

	log.Println("The database opens successfully")
	return db
}

func CloseDatabase(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("problem getting sql.DB from gorm.DB: %v", err)
		return
	}

	err = sqlDB.Close()
	if err != nil {
		log.Printf("problem closing database: %v", err)
	} else {
		log.Println("Database closed successfully")
	}
}

func migrateDB(db *gorm.DB) error {
	log.Println("Starting database migration...")

	err := db.AutoMigrate(&Account{}, &Category{}, &SubCategory{}, &Transaction{}, &TransferTransaction{})
	if err != nil {
		return fmt.Errorf("failed to create database: %v", err)
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

	log.Printf("Successfully got %d accounts from db", len(accounts))
	return accounts, nil
}

func GetCategories(db *gorm.DB) ([]Category, error) {
	log.Println("Getting all categories from database")

	var caregories []Category
	result := db.Find(&caregories)

	if result.Error != nil {
		log.Printf("Error getting accounts: %v", result.Error)
		return nil, fmt.Errorf("problem with get categories in db: %v", result.Error)
	}

	log.Printf("Successfully got %d categories from db", len(caregories))
	return caregories, nil
}

func GetSubCategories(db *gorm.DB) ([]SubCategory, error) {
	log.Println("Getting all sub_categories from database")

	var sub_categories []SubCategory
	result := db.Find(&sub_categories)

	if result.Error != nil {
		log.Printf("Error getting accounts: %v", result.Error)
		return nil, fmt.Errorf("problem with get sub_categories in db: %v", result.Error)
	}

	log.Printf("Successfully got %d sub_categories from db", len(sub_categories))
	return sub_categories, nil
}

func GetTransactions(db *gorm.DB) ([]Transaction, error) {
	log.Println("Getting all transactions from database")

	var transactions []Transaction
	result := db.
		Preload("Account").
		Preload("Category").
		Preload("SubCategory").
		Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}

	log.Printf("Successfully got %d sub_categories from db", len(transactions))
	return transactions, nil
}
