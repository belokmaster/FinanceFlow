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

func DeleteTransfer(db *gorm.DB, id int) error {
	log.Printf("Deleting transfer with ID: %d", id)

	return db.Transaction(func(txDB *gorm.DB) error {
		var tx TransferTransaction
		err := txDB.First(&tx, id).Error
		if err != nil {
			return fmt.Errorf("transfer with ID %d not found", id)
		}

		var fromAccount Account
		err = txDB.First(&fromAccount, tx.AccountID).Error
		if err != nil {
			return fmt.Errorf("from account with ID %d not found", tx.AccountID)
		}

		var toAccount Account
		err = txDB.First(&toAccount, tx.TransferAccountID).Error
		if err != nil {
			return fmt.Errorf("to account with ID %d not found", tx.TransferAccountID)
		}

		log.Printf("Found transfer: From=%s, To=%s, Amount=%.2f", fromAccount.Name, toAccount.Name, tx.Amount)
		log.Printf("Old balances: From=%.2f, To=%.2f", fromAccount.Balance, toAccount.Balance)

		fromAccount.Balance += tx.Amount
		toAccount.Balance -= tx.Amount

		log.Printf("New balances: From=%.2f, To=%.2f", fromAccount.Balance, toAccount.Balance)

		if err := txDB.Save(&fromAccount).Error; err != nil {
			return fmt.Errorf("failed to update from account balance: %v", err)
		}

		if err := txDB.Save(&toAccount).Error; err != nil {
			return fmt.Errorf("failed to update to account balance: %v", err)
		}

		if err := txDB.Delete(&TransferTransaction{}, id).Error; err != nil {
			return fmt.Errorf("failed to delete transfer: %v", err)
		}

		log.Printf("Transfer deleted and balances adjusted: From=%s, To=%s", fromAccount.Name, toAccount.Name)
		return nil
	})
}

func UpdateTransfer(db *gorm.DB, newTx TransferTransaction) error {
	log.Printf("Updating transfer with ID: %d", newTx.ID)

	return db.Transaction(func(txDB *gorm.DB) error {
		var oldTx TransferTransaction
		err := txDB.First(&oldTx, newTx.ID).Error
		if err != nil {
			return fmt.Errorf("transfer not found: %v", err)
		}

		if newTx.AccountID == newTx.TransferAccountID {
			return fmt.Errorf("cannot transfer to the same account")
		}

		var oldFromAccount Account
		err = txDB.First(&oldFromAccount, oldTx.AccountID).Error
		if err != nil {
			return fmt.Errorf("old from account not found: %v", err)
		}

		var oldToAccount Account
		err = txDB.First(&oldToAccount, oldTx.TransferAccountID).Error
		if err != nil {
			return fmt.Errorf("old to account not found: %v", err)
		}

		var newFromAccount Account
		err = txDB.First(&newFromAccount, newTx.AccountID).Error
		if err != nil {
			return fmt.Errorf("new from account not found: %v", err)
		}

		var newToAccount Account
		err = txDB.First(&newToAccount, newTx.TransferAccountID).Error
		if err != nil {
			return fmt.Errorf("new to account not found: %v", err)
		}

		log.Printf("Old transfer: From=%s, To=%s, Amount=%.2f", oldFromAccount.Name, oldToAccount.Name, oldTx.Amount)
		log.Printf("New transfer: From=%s, To=%s, Amount=%.2f", newFromAccount.Name, newToAccount.Name, newTx.Amount)

		oldFromAccount.Balance += oldTx.Amount
		oldToAccount.Balance -= oldTx.Amount

		if newFromAccount.Balance < newTx.Amount {
			return fmt.Errorf("insufficient funds in new from account: available %.2f, required %.2f", newFromAccount.Balance, newTx.Amount)
		}

		newFromAccount.Balance -= newTx.Amount
		newToAccount.Balance += newTx.Amount

		accountsToSave := []*Account{&oldFromAccount, &oldToAccount, &newFromAccount, &newToAccount}
		for _, acc := range accountsToSave {
			if err := txDB.Save(acc).Error; err != nil {
				return fmt.Errorf("failed to update account %s balance: %v", acc.Name, err)
			}
		}

		err = txDB.Model(&TransferTransaction{}).Where("id = ?", newTx.ID).Updates(newTx).Error
		if err != nil {
			return fmt.Errorf("failed to update transfer: %v", err)
		}

		log.Printf("Transfer updated successfully: ID=%d", newTx.ID)
		return nil
	})
}

func GetTransfer(db *gorm.DB, id int) (TransferTransaction, error) {
	var transfer TransferTransaction
	err := db.Preload("Account").Preload("TransferAccount").First(&transfer, id).Error
	if err != nil {
		return TransferTransaction{}, err
	}
	return transfer, nil
}

func GetTransfers(db *gorm.DB) ([]TransferTransaction, error) {
	var transfers []TransferTransaction
	err := db.Preload("Account").Preload("TransferAccount").Order("date DESC").Find(&transfers).Error
	if err != nil {
		return nil, err
	}
	return transfers, nil
}
