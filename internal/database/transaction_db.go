package database

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

func AddTransaction(db *gorm.DB, tx Transaction) error {
	log.Printf("Adding new transaction: AccountID=%d, Type=%v, Amount=%.2f, Comment=%s",
		tx.AccountID, tx.Type, tx.Amount, tx.Comment)

	return db.Transaction(func(txDB *gorm.DB) error {
		// add record to transactions
		err := txDB.Create(&tx).Error
		if err != nil {
			log.Printf("Error creating transaction record: %v", err)
			return err
		}

		log.Printf("Transaction record created successfully: ID=%d", tx.ID)

		// parse account info
		var account Account
		err = txDB.First(&account, tx.AccountID).Error
		if err != nil {
			log.Printf("Error finding account ID %d: %v", tx.AccountID, err)
			return err
		}

		log.Printf("Found account: ID=%d, Name=%s, Old Balance=%.2f", account.ID, account.Name, account.Balance)

		if tx.Amount < 0 {
			return fmt.Errorf("transaction amount must be positive")
		}

		switch tx.Type {
		case Income:
			account.Balance += tx.Amount
			log.Printf("Income transaction: adding %.2f to balance", tx.Amount)
		case Expense:
			account.Balance -= tx.Amount
			log.Printf("Expense transaction: subtracting %.2f from balance", tx.Amount)
		default:
			log.Printf("Error: unknown transaction type: %v", tx.Type)
			return fmt.Errorf("unknown transaction type: %v", tx.Type)
		}

		log.Printf("New account balance: %.2f", account.Balance)
		err = txDB.Save(&account).Error
		if err != nil {
			log.Printf("Error updating account balance: %v", err)
			return err
		}

		log.Printf("Account balance updated successfully for account ID %d", account.ID)
		log.Printf("Transaction completed successfully: ID=%d", tx.ID)
		return nil
	})
}

func DeleteTransaction(db *gorm.DB, id int) error {
	log.Printf("Deleting transaction with ID: %d", id)

	return db.Transaction(func(txDB *gorm.DB) error {
		var tx Transaction
		err := txDB.First(&tx, id).Error
		if err != nil {
			return fmt.Errorf("transaction with ID %d not found", id)
		}

		var acc Account
		err = txDB.First(&acc, tx.AccountID).Error
		if err != nil {
			return fmt.Errorf("account with ID %d not found", tx.AccountID)
		}

		// return past balance
		switch tx.Type {
		case Income:
			acc.Balance -= tx.Amount
		case Expense:
			acc.Balance += tx.Amount
		}

		if err := txDB.Save(&acc).Error; err != nil {
			return fmt.Errorf("failed to update balance: %v", err)
		}

		if err := txDB.Delete(&Transaction{}, id).Error; err != nil {
			return fmt.Errorf("failed to delete transaction: %v", err)
		}

		log.Printf("Transaction deleted and balance adjusted for account %d", acc.ID)
		return nil
	})
}

func UpdateTransaction(db *gorm.DB, newTx Transaction) error {
	return db.Transaction(func(txDB *gorm.DB) error {
		var oldTx Transaction
		err := txDB.First(&oldTx, newTx.ID).Error
		if err != nil {
			return fmt.Errorf("transaction not found: %v", err)
		}

		var acc Account
		err = txDB.First(&acc, oldTx.AccountID).Error
		if err != nil {
			return fmt.Errorf("account not found: %v", err)
		}

		// back the old change
		if oldTx.Type == Income {
			acc.Balance -= oldTx.Amount
		} else {
			acc.Balance += oldTx.Amount
		}

		// approve new balance
		if newTx.Type == Income {
			acc.Balance += newTx.Amount
		} else {
			acc.Balance -= newTx.Amount
		}

		err = txDB.Save(&acc).Error
		if err != nil {
			return fmt.Errorf("failed to update account balance: %v", err)
		}

		err = txDB.Model(&Transaction{}).Where("id = ?", newTx.ID).Updates(newTx).Error
		if err != nil {
			return fmt.Errorf("failed to update transaction: %v", err)
		}

		return nil
	})
}
