package account

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountRepository interface {
	CreateAccount(accountData Account) (*Account, error)
	FindAccountByUserId(userId string) (*Account, error)
	FindAccountByAccountNumber(accountNumber string) (*Account, error)
	UpdateBalance(accountData Account) (*Account, error)
}

type Account struct {
	gorm.Model
	ID            string  `gorm:"primary_key"`
	UserID        string  `gorm:"user_id"`
	AccountNumber string  `gorm:"account_number"`
	Balance       float64 `gorm:"balance"`
}

func (a *Account) BeforeCreate(db *gorm.DB) (err error) {
	a.ID = uuid.New().String()
	//a.AccountNumber = shared.GenerateAccountNumber()
	return err
}
