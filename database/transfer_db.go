package database

import (
	"fmt"

	"gorm.io/gorm"
)

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
