package service

import (
	"github.com/FaisalMashuri/my-wallet/internal/domain/account"
	"github.com/FaisalMashuri/my-wallet/shared"
	"github.com/FaisalMashuri/my-wallet/shared/contract"
	"github.com/pkg/errors"
)

type accountServiceImpl struct {
	repo account.AccountRepository
}

func NewService(repo account.AccountRepository) account.AccountService {
	return &accountServiceImpl{
		repo: repo,
	}
}

func (s *accountServiceImpl) CreateAccount(userId string) (*account.Account, error) {
	totalData, err := s.repo.CountAccountNumberByUserId(userId)
	if totalData >= 2 {
		return nil, errors.New(contract.ErrLimitAccountOpen)
	}
	accountData := account.Account{
		UserID:        userId,
		Balance:       0,
		AccountNumber: shared.GenerateAccountNumber(),
	}
	accountRes, err := s.repo.CreateAccount(accountData)
	if err != nil {
		return nil, errors.New(contract.ErrInternalServer)
	}
	return accountRes, nil
}
