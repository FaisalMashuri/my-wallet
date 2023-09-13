package repository

import (
	"github.com/FaisalMashuri/my-wallet/internal/domain/account"
	"gorm.io/gorm"
)

type accountRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) account.AccountRepository {
	return &accountRepository{
		db: db,
	}
}

func (a *accountRepository) CreateAccount(accountData account.Account) (*account.Account, error) {
	//TODO implement me
	err := a.db.Debug().Create(&accountData).Error
	if err != nil {
		return nil, err
	}
	return &accountData, nil
}

func (a *accountRepository) FindAccountByUserId(userId string) (*account.Account, error) {
	//TODO implement me
	var accountModel account.Account
	err := a.db.Debug().Take(&accountModel, "user_id = ?", userId).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &accountModel, nil

}

func (a *accountRepository) FindAccountByAccountNumber(accountNumber string) (*account.Account, error) {
	//TODO implement me
	var accountModel account.Account
	err := a.db.Debug().Take(&accountModel, "account_number = ?", accountNumber).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &accountModel, nil
}

func (a *accountRepository) FindAllAccountsByUserId(userId string) (accountModel []*account.Account, err error) {
	err = a.db.Debug().Find(&accountModel, "user_id = ?", userId).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return accountModel, nil
		}
		return nil, err
	}
	return
}

func (a *accountRepository) UpdateBalance(accountData account.Account) (*account.Account, error) {
	//TODO implement me
	err := a.db.Debug().Updates(&accountData).Error
	if err != nil {
		return nil, err
	}
	return &accountData, nil
}

func (a *accountRepository) CountAccountNumberByUserId(userId string) (totalData int64, err error) {
	err = a.db.Debug().Model(&account.Account{}).Where("user_id = ?", userId).Count(&totalData).Error
	if err != nil {
		return 0, err
	}
	return totalData, nil
}
