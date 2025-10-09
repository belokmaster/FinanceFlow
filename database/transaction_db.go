package database

import (
	"fmt"

	"gorm.io/gorm"
)

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
