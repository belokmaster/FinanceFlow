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
	Transfer
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

type TransferTransaction struct {
	ID                uint `gorm:"primaryKey;autoIncrement"`
	AccountID         uint
	Account           Account
	TransferAccountID uint
	TransferAccount   Account
	Type              TypeTransaction
	Amount            float64
	Comment           string
	Date              time.Time
}

func InitDB(db *gorm.DB) error {
	err := db.AutoMigrate(&Account{}, &Category{}, &SubCategory{}, &Transaction{}, &TransferTransaction{})
	if err != nil {
		log.Fatalf("problem to create gorm bd: %v", err)
	}

	log.Println("The database created successfully")
	return nil
}

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

func TransferToAnotherAcc(db *gorm.DB, tx TransferTransaction) error {
	if tx.Type != Transfer {
		return fmt.Errorf("problem with transaction types")
	}

	if tx.AccountID == tx.TransferAccountID {
		return fmt.Errorf("transfer to yourself")
	}

	if tx.Amount < 0 {
		return fmt.Errorf("transfer is negative")
	}

	return db.Transaction(func(txDB *gorm.DB) error {
		var fromAccount Account
		var toAccount Account

		err := txDB.First(&fromAccount, tx.AccountID).Error
		if err != nil {
			return err
		}

		err = txDB.First(&toAccount, tx.TransferAccountID).Error
		if err != nil {
			return err
		}

		if fromAccount.Balance < tx.Amount {
			return fmt.Errorf("balance less then transfer sum")
		}

		// block change balance each accounts
		fromAccount.Balance -= tx.Amount
		err = txDB.Save(&fromAccount).Error
		if err != nil {
			return err
		}

		toAccount.Balance += tx.Amount
		err = txDB.Save(&toAccount).Error
		if err != nil {
			return err
		}

		// create the record after all checks upper
		err = txDB.Create(&tx).Error
		if err != nil {
			return err
		}

		return nil
	})
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
