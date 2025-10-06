package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

type Account struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	Name    string `gorm:"uniqueIndex;not null"`
	Balance float64
}

type Category struct {
	ID   uint   `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"uniqueIndex;not null"`
}

type SubCategory struct {
	ID         uint   `gorm:"primaryKey;autoIncrement"`
	Name       string `gorm:"uniqueIndex;not null"`
	CategoryID uint
	Category   Category `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type TypeTransaction int

const (
	Income TypeTransaction = iota
	Expense
)

type Transaction struct {
	ID            uint `gorm:"primaryKey;autoIncrement"`
	AccountID     uint
	Account       Account
	CategoryID    uint
	Category      Category
	SubCategoryID uint
	SubCategory   SubCategory
	Type          TypeTransaction
	Amount        float64
	Comment       string
	Date          time.Time
}

func InitDB(db *gorm.DB) error {
	err := db.AutoMigrate(&Account{}, &Category{}, &SubCategory{}, &Transaction{})
	if err != nil {
		log.Fatalf("problem to create gorm bd: %v", err)
	}

	log.Println("The database created successfully")
	return nil
}

func AddTransaction(db *gorm.DB, tx Transaction) error {
	return db.Transaction(func(txDB *gorm.DB) error {
		// add record to transactions
		err := txDB.Create(&tx).Error
		if err != nil {
			return err
		}

		// parse account info
		var account Account
		err = txDB.First(&account, tx.AccountID).Error
		if err != nil {
			return err
		}

		switch tx.Type {
		case Income:
			account.Balance += tx.Amount
		case Expense:
			account.Balance -= tx.Amount
		default:
			return fmt.Errorf("unknown transaction type: %v", tx.Type)
		}

		err = txDB.Save(&account).Error
		if err != nil {
			return err
		}

		return nil
	})
}
