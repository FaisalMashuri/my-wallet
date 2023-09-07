package transaction

import (
	"github.com/FaisalMashuri/my-wallet/internal/domain/transaction/dto/request"
	"github.com/FaisalMashuri/my-wallet/internal/domain/transaction/dto/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type TransactionRepository interface {
	CreateTransaction(transactionModel Transaction) (*Transaction, error)
	CreateTransactionInquiry(inquiry TransactionInquiry) (*TransactionInquiry, error)
	FindTransactionInquiry(inquiryKey string) (*TransactionInquiry, error)
	DeleteInquiry(inquiryKey string) error
}

type TransactionService interface {
	TranferInquiry(InquiryReq request.TransferInquiryReq, ctx *fiber.Ctx) (*response.TransferInquiryRes, error)
	TransferInquiryExec(InquiryExecReq request.TransferInquiryExec, ctx *fiber.Ctx) error
}

type TransactionController interface {
	TransferInquiry(ctx *fiber.Ctx) error
	TransferExec(ctx *fiber.Ctx) error
}

type TransactionInquiry struct {
	gorm.Model
	ID         string    `gorm:"primary_key"`
	InquiryKey string    `gorm:"inquiry_key"`
	Value      string    `gorm:"value;type:jsonb"`
	usedAt     time.Time `gorm:"used_at"`
}

type Transaction struct {
	gorm.Model
	ID              string    `gorm:"primary_key"`
	AccountID       string    `gorm:"account_id"`
	Amount          float64   `gorm:"amount"`
	SofNumber       string    `gorm:"sof_number"`
	DofNumber       string    `gorm:"dof_number"`
	TransactionType string    `gorm:"transaction_type"`
	TransactionTime time.Time `gorm:"transaction_time"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New().String()
	t.TransactionTime = time.Now()
	return
}

func (t *TransactionInquiry) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New().String()
	return
}
