package repository

import (
	"github.com/FaisalMashuri/my-wallet/internal/domain/mpin"
	"github.com/FaisalMashuri/my-wallet/shared/contract"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type mPinRepostoryImpl struct {
	db *gorm.DB
}

func (m *mPinRepostoryImpl) CreatePin(modelPin mpin.Pin) (err error) {
	//TODO implement me
	err = m.db.Debug().Create(&modelPin).Error
	if err != nil {
		return err
	}
	return
}

func (m *mPinRepostoryImpl) FindByUserId(userId string) (pin mpin.Pin, err error) {
	//TODO implement me
	err = m.db.Debug().First(&pin, "user_id = ?", userId).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return pin, errors.New(contract.ErrRecordNotFound)
		}
		return mpin.Pin{}, errors.New(contract.ErrInternalServer)
	}
	return
}

func NewRepository(db *gorm.DB) mpin.PinRepository {
	return &mPinRepostoryImpl{
		db: db,
	}
}
