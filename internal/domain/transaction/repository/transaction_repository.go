package repository

import (
	"github.com/FaisalMashuri/my-wallet/internal/domain/transaction"
	"github.com/FaisalMashuri/my-wallet/shared/contract"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

func (t *transactionRepository) DeleteInquiry(inquiryKey string) error {
	//TODO implement me
	var inquiryModel transaction.TransactionInquiry
	err := t.db.Debug().Delete(&inquiryModel, "inquiry_key = ?", inquiryKey).Error
	if err != nil {
		return err
	}
	return nil

}

func (t transactionRepository) FindTransactionInquiry(inquiryKey string) (*transaction.TransactionInquiry, error) {
	//TODO implement me
	var inquiry transaction.TransactionInquiry
	err := t.db.Debug().First(&inquiry, "inquiry_key = ?", inquiryKey).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.New(contract.ErrInternalServer)
	}
	return &inquiry, nil
}

func (t transactionRepository) CreateTransaction(transactionModel transaction.Transaction) (*transaction.Transaction, error) {
	//TODO implement me
	err := t.db.Debug().Create(&transactionModel).Error
	if err != nil {
		return nil, err
	}
	return &transactionModel, nil

}

func (t transactionRepository) CreateTransactionInquiry(inquiry transaction.TransactionInquiry) (*transaction.TransactionInquiry, error) {
	//TODO implement me
	err := t.db.Debug().Create(&inquiry).Error
	if err != nil {
		return nil, err
	}
	return &inquiry, nil
}

func NewRepository(db *gorm.DB) transaction.TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}
