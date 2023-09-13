package account

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountRepository interface {
	CreateAccount(accountData Account) (*Account, error)
	FindAccountByUserId(userId string) (*Account, error)
	FindAllAccountsByUserId(userId string) ([]*Account, error)
	FindAccountByAccountNumber(accountNumber string) (*Account, error)
	UpdateBalance(accountData Account) (*Account, error)
	CountAccountNumberByUserId(userId string) (int64, error)
}

type AccountService interface {
	CreateAccount(userId string) (*Account, error)
}
type AccountController interface {
	CreateAccount(ctx *fiber.Ctx) error
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
