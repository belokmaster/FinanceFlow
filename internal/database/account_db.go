package database

import (
	"fmt"
	"log"
	"strings"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func CreateNewAccount(db *gorm.DB, acc Account) error {
	if strings.TrimSpace(acc.Name) == "" {
		err := fmt.Errorf("account name cannot be empty")
		log.Println("Error creating account:", err)
		return err
	}

	log.Printf("Creating new account: %s", acc.Name)

	result := db.Create(&acc)
	if result.Error != nil {
		log.Printf("Error creating account %s: %v", acc.Name, result.Error)
		return fmt.Errorf("problem to create a new acc: %v", result.Error)
	}

	log.Printf("Account %s created successfully with ID=%d", acc.Name, acc.ID)
	return nil
}

func DeleteAccount(db *gorm.DB, id int) error {
	log.Printf("Deleting account with ID: %d", id)

	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("Error starting transaction: %v", tx.Error)
		return fmt.Errorf("problem starting transaction: %v", tx.Error)
	}

	// delete transactions and transfers before delete account
	// it needs to recover balance
	var transactions []Transaction
	if err := tx.Where("account_id = ?", id).Find(&transactions).Error; err != nil {
		tx.Rollback()
		log.Printf("Error finding transactions for account ID %d: %v", id, err)
		return fmt.Errorf("problem finding transactions: %v", err)
	}

	for _, transaction := range transactions {
		if err := DeleteTransaction(tx, int(transaction.ID)); err != nil {
			tx.Rollback()
			log.Printf("Error deleting transaction ID %d for account ID %d: %v", transaction.ID, id, err)
			return fmt.Errorf("problem deleting transaction %d: %v", transaction.ID, err)
		}
	}
	log.Printf("Deleted %d transactions for account ID %d", len(transactions), id)

	var transfers []TransferTransaction
	if err := tx.Where("account_id = ? OR transfer_account_id = ?", id, id).Find(&transfers).Error; err != nil {
		tx.Rollback()
		log.Printf("Error finding transfers for account ID %d: %v", id, err)
		return fmt.Errorf("problem finding transfers: %v", err)
	}
	for _, tr := range transfers {
		if err := DeleteTransfer(tx, int(tr.ID)); err != nil {
			tx.Rollback()
			log.Printf("Error deleting transfer ID %d for account ID %d: %v", tr.ID, id, err)
			return fmt.Errorf("problem deleting transfer %d: %v", tr.ID, err)
		}
	}
	log.Printf("Deleted %d transfers for account ID %d", len(transfers), id)

	result := tx.Delete(&Account{}, id)
	if result.Error != nil {
		tx.Rollback()
		log.Printf("Error deleting account ID %d: %v", id, result.Error)
		return fmt.Errorf("problem with delete account in db: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		log.Printf("Account with ID %d not found for deletion", id)
		return fmt.Errorf("account with id %d not found", id)
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction for account ID %d: %v", id, err)
		return fmt.Errorf("problem committing transaction: %v", err)
	}

	log.Printf("Account ID %d and all related transactions/transfers deleted successfully", id)
	return nil
}

func ChangeAccountColor(db *gorm.DB, id int, newColor string) error {
	log.Printf("Changing color for account ID %d to: %s", id, newColor)

	if newColor == "" {
		log.Printf("Error: empty color for account ID %d", id)
		return fmt.Errorf("empty color")
	}

	newColor = strings.TrimPrefix(newColor, "#")

	result := db.Model(&Account{}).Where("id = ?", id).Update("color", newColor)
	if result.Error != nil {
		log.Printf("Error updating color for account ID %d: %v", id, result.Error)
		return fmt.Errorf("problem with change color in db: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		log.Printf("No rows affected when updating color for account ID %d", id)
		return fmt.Errorf("no rows were updated - account may not exist")
	}

	log.Printf("Color updated successfully for account ID %d", id)
	return nil
}

func ChangeAccountIcon(db *gorm.DB, id int, icon_id int) error {
	log.Printf("Changing icon for account ID %d to icon ID: %d", id, icon_id)

	if _, ok := IconAccountFiles[TypeIcons(icon_id)]; !ok {
		log.Printf("Error: icon ID %d does not exist", icon_id)
		return fmt.Errorf("icon with ID %d does not exist", icon_id)
	}

	var acc Account
	err := db.First(&acc, id).Error
	if err != nil {
		log.Printf("Error finding account ID %d for icon change: %v", id, err)
		return fmt.Errorf("account with ID %d not found", id)
	}

	log.Printf("Found account: ID=%d, Name=%s, Old icon=%d, New icon=%d", acc.ID, acc.Name, acc.IconCode, icon_id)

	result := db.Model(&Account{}).Where("id = ?", id).Update("icon_code", TypeIcons(icon_id))
	if result.Error != nil {
		log.Printf("Error updating icon for account ID %d: %v", id, result.Error)
		return fmt.Errorf("problem with changing icon in db: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		log.Printf("No rows affected when updating icon for account ID %d", id)
		return fmt.Errorf("no rows were updated - account may not exist")
	}

	log.Printf("Icon updated successfully for account ID %d", id)
	return nil
}

func ChangeAccountName(db *gorm.DB, id int, newName string) error {
	var acc Account
	err := db.First(&acc, id).Error
	if err != nil {
		return fmt.Errorf("account with ID %d not found", id)
	}

	result := db.Model(&Account{}).Where("id = ?", id).Update("Name", newName)
	if result.Error != nil {
		return fmt.Errorf("problem with change name in db: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		log.Printf("No rows affected when updating name for account ID %d", id)
		return fmt.Errorf("no rows were updated - account may not exist")
	}

	return nil
}

func ChangeAccountBalance(db *gorm.DB, id int, newBalance float64) error {
	var acc Account
	err := db.First(&acc, id).Error
	if err != nil {
		return fmt.Errorf("account with ID %d not found", id)
	}

	balance := decimal.NewFromFloat(newBalance)
	result := db.Model(&Account{}).Where("id = ?", id).Update("Balance", balance)
	if result.Error != nil {
		return fmt.Errorf("problem with change balance in db: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		log.Printf("No rows affected when updating balance for account ID %d", id)
		return fmt.Errorf("no rows were updated - account may not exist")
	}

	return nil
}
