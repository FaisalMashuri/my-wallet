package repository

import (
	"github.com/FaisalMashuri/my-wallet/internal/domain/topup"
	"gorm.io/gorm"
)

type topUpRepository struct {
	db *gorm.DB
}

func (t *topUpRepository) FindById(id string) (topUp *topup.TopUp, err error) {
	//TODO implement me
	err = t.db.Debug().First(&topUp).Error
	if err != nil {
		return nil, err
	}
	return
}

func (t *topUpRepository) Insert(topUp *topup.TopUp) error {
	//TODO implement me
	err := t.db.Debug().Create(&topUp).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *topUpRepository) Update(topUp *topup.TopUp) error {
	//TODO implement me
	err := t.db.Debug().Updates(&topUp).Error
	if err != nil {
		return err
	}
	return nil
}

func NewRepository(db *gorm.DB) topup.TopUpRepository {
	return &topUpRepository{
		db: db,
	}
}
