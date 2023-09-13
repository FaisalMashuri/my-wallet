package service

import (
	"fmt"
	"github.com/FaisalMashuri/my-wallet/internal/domain/mpin"
	"github.com/FaisalMashuri/my-wallet/internal/domain/mpin/dto/request"
	"github.com/FaisalMashuri/my-wallet/shared/contract"
	"github.com/pkg/errors"
)

type mPinServiceImpl struct {
	repo mpin.PinRepository
}

func (m mPinServiceImpl) CreatePin(request request.PinRequest) error {
	//TODO implement me
	pinModel := mpin.Pin{
		UserID: request.UserId,
		Pin:    request.Pin,
	}
	err := m.repo.CreatePin(pinModel)
	if err != nil {
		return errors.New(contract.ErrInternalServer)
	}
	return nil
}

func (m mPinServiceImpl) ValidatePin(request *request.ValidatePinReq) error {
	//TODO implement me
	fmt.Println("Vaildasi start")
	dataPin, err := m.repo.FindByUserId(request.UserId)
	if dataPin == (mpin.Pin{}) {
		if err != nil {
			return err
		}
	}
	if dataPin.Pin != request.Pin {
		return errors.New(contract.ErrInvalidPin)
	}
	return err
}

func NewService(repo mpin.PinRepository) mpin.PinService {
	return &mPinServiceImpl{
		repo: repo,
	}
}
