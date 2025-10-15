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
	dbExists := true
	if err != nil {
		if os.IsNotExist(err) {
			dbExists = false
		} else {
			log.Fatalf("problem with db file: %v", err)
		}
	}

	// open the db using gorm
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("problem connecting to db: %v", err)
	}

	log.Println("Running migrations...")

	err = migrateDB(db)
	if err != nil {
		log.Fatalf("failed to migrate the db: %v", err)
	}

	if !dbExists {
		log.Println("Database created successfully")
	} else {
		log.Println("Database opened successfully")
	}

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

	log.Println("Database migration completed successfully")
	return nil
}

func GetAccounts(db *gorm.DB) ([]Account, error) {
	log.Println("Getting all accounts from database")

	var accounts []Account
	result := db.Find(&accounts)

	if result.Error != nil {
		return nil, fmt.Errorf("problem with get accounts in db: %v", result.Error)
	}

	log.Printf("Successfully got %d accounts from db", len(accounts))
	return accounts, nil
}

func GetCategories(db *gorm.DB) ([]Category, error) {
	log.Println("Getting all categories from database")

	var categories []Category
	result := db.Find(&categories)

	if result.Error != nil {
		return nil, fmt.Errorf("problem with get categories in db: %v", result.Error)
	}

	log.Printf("Successfully got %d categories from db", len(categories))
	return categories, nil
}

func GetSubCategories(db *gorm.DB) ([]SubCategory, error) {
	log.Println("Getting all sub_categories from database")

	var subCategories []SubCategory
	result := db.Find(&subCategories)

	if result.Error != nil {
		return nil, fmt.Errorf("problem with get sub_categories in db: %v", result.Error)
	}

	log.Printf("Successfully got %d sub_categories from db", len(subCategories))
	return subCategories, nil
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

	log.Printf("Successfully got %d transactions from db", len(transactions))
	return transactions, nil
}
