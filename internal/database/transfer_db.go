package database

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

func TransferToAnotherAcc(db *gorm.DB, tx TransferTransaction) error {
	log.Printf("Starting transfer: Amount=%.2f, From Account=%d, To Account=%d", tx.Amount, tx.AccountID, tx.TransferAccountID)

	if tx.Type != Transfer {
		log.Printf("Error: invalid transaction type for transfer: %v", tx.Type)
		return fmt.Errorf("problem with transaction types")
	}

	if tx.AccountID == tx.TransferAccountID {
		log.Printf("Error: transfer to same account ID %d", tx.AccountID)
		return fmt.Errorf("transfer to yourself")
	}

	if tx.Amount < 0 {
		log.Printf("Error: negative transfer amount: %.2f", tx.Amount)
		return fmt.Errorf("transfer is negative")
	}

	return db.Transaction(func(txDB *gorm.DB) error {
		var fromAccount Account
		var toAccount Account

		log.Printf("Fetching source account ID: %d", tx.AccountID)
		err := txDB.First(&fromAccount, tx.AccountID).Error
		if err != nil {
			log.Printf("Error finding source account ID %d: %v", tx.AccountID, err)
			return err
		}

		log.Printf("Fetching destination account ID: %d", tx.TransferAccountID)

		err = txDB.First(&toAccount, tx.TransferAccountID).Error
		if err != nil {
			log.Printf("Error finding destination account ID %d: %v", tx.TransferAccountID, err)
			return err
		}

		log.Printf("Source account: ID=%d, Name=%s, Balance=%.2f", fromAccount.ID, fromAccount.Name, fromAccount.Balance)
		log.Printf("Destination account: ID=%d, Name=%s, Balance=%.2f", toAccount.ID, toAccount.Name, toAccount.Balance)

		if fromAccount.Balance < tx.Amount {
			log.Printf("Error: insufficient funds. Required: %.2f, Available: %.2f", tx.Amount, fromAccount.Balance)
			return fmt.Errorf("balance less then transfer sum")
		}

		// block change balance each accounts
		log.Printf("Updating source account balance: %.2f -> %.2f", fromAccount.Balance, fromAccount.Balance-tx.Amount)

		fromAccount.Balance -= tx.Amount
		err = txDB.Save(&fromAccount).Error
		if err != nil {
			log.Printf("Error updating source account balance: %v", err)
			return err
		}

		log.Printf("Updating destination account balance: %.2f -> %.2f", toAccount.Balance, toAccount.Balance+tx.Amount)

		toAccount.Balance += tx.Amount
		err = txDB.Save(&toAccount).Error
		if err != nil {
			log.Printf("Error updating destination account balance: %v", err)
			return err
		}

		// create the record after all checks upper
		log.Printf("Creating transfer transaction record")

		err = txDB.Create(&tx).Error
		if err != nil {
			log.Printf("Error creating transfer transaction record: %v", err)
			return err
		}

		log.Printf("Transfer completed successfully: ID=%d, From Account=%s, To Account=%s, Amount=%.2f", tx.ID, fromAccount.Name, toAccount.Name, tx.Amount)
		log.Printf("Final balances. Source: %.2f, Destination: %.2f", fromAccount.Balance, toAccount.Balance)

		return nil
	})
}
