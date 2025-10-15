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

	result := db.Delete(&Transaction{}, id)
	if result.Error != nil {
		log.Printf("Error deleting transaction ID %d: %v", id, result.Error)
		return fmt.Errorf("problem with delete transaction in db: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		log.Printf("Transaction with ID %d not found for deletion", id)
		return fmt.Errorf("Transaction with id %d not found", id)
	}

	log.Printf("Transaction deleted successfully: ID=%d", id)
	return nil
}

func UpdateTransaction(db *gorm.DB, transaction Transaction) error {
	result := db.Model(&Transaction{}).Where("id = ?", transaction.ID).Updates(map[string]interface{}{
		"account_id":      transaction.AccountID,
		"category_id":     transaction.CategoryID,
		"sub_category_id": transaction.SubCategoryID,
		"type":            transaction.Type,
		"amount":          transaction.Amount,
		"comment":         transaction.Comment,
		"date":            transaction.Date,
	})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("transaction not found")
	}

	return nil
}
