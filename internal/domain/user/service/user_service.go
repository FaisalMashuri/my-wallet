package service

import (
	"fmt"
	"github.com/FaisalMashuri/my-wallet/internal/domain/account"
	dtoResAccount "github.com/FaisalMashuri/my-wallet/internal/domain/account/dto/response"
	"github.com/FaisalMashuri/my-wallet/internal/domain/user"
	"github.com/FaisalMashuri/my-wallet/internal/domain/user/dto/request"
	"github.com/FaisalMashuri/my-wallet/internal/domain/user/dto/response"

	"github.com/FaisalMashuri/my-wallet/shared"
	"github.com/FaisalMashuri/my-wallet/shared/contract"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo        user.UserRepository
	repoAccount account.AccountRepository
	log         *logrus.Logger
}

func NewService(repo user.UserRepository, log *logrus.Logger, repoAccount account.AccountRepository) user.UserService {
	return &userService{
		repo:        repo,
		repoAccount: repoAccount,
		log:         log,
	}
}

func (s *userService) Login(userRequest *request.LoginRequest) (*response.LoginResponse, error) {
	userData, err := s.repo.FindUserByEmail(userRequest.Email)
	if userData == nil {
		if err != nil {
			s.log.WithField("error", err.Error()).Info("failed to find user")
			return nil, errors.New(contract.ErrInternalServer)
		}
		s.log.Info("User not found")
		return nil, errors.New(contract.ErrRecordNotFound)
	}

	err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(userRequest.Password))
	if err != nil {
		s.log.Info("Password mismatch")
		return nil, errors.New(contract.ErrPasswordNotMatch)
	}

	accessToken, err := shared.GenerateAccessToken(userData)
	if err != nil {
		s.log.Info("faile generate token")
		return nil, errors.New(contract.ErrUnexpectedError)
	}
	refeshToken, err := shared.GenerateRefreshToken(userData)
	if err != nil {
		s.log.Info("faile generate token")
		return nil, errors.New(contract.ErrUnexpectedError)
	}

	res := response.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refeshToken,
	}
	return &res, nil
}

func (s *userService) RegisterUser(userRequest *request.RegisterRequest) (userData *user.User, err error) {
	userData, err = s.repo.FindUserByEmail(userRequest.Email)
	if userData != nil {
		if err != nil {
			s.log.WithField("error", err.Error()).Info("Error find email")
			return nil, errors.New(contract.ErrInternalServer)

		}
		s.log.Info("Email already registered")
		return nil, errors.New(contract.ErrEmailAlreadyRegister)
	}

	userModel := userData.FromRegistRequest(userRequest)
	userData, err = s.repo.CreateUser(&userModel)
	if err != nil {
		s.log.WithField("error", err.Error()).Info("failed to create user")
		return nil, err
	}

	accountData := account.Account{
		UserID:        userData.ID,
		Balance:       0,
		AccountNumber: shared.GenerateAccountNumber(),
	}

	_, err = s.repoAccount.CreateAccount(accountData)
	if err != nil {
		s.log.Info("failed to create account")
		return nil, errors.New(contract.ErrInternalServer)
	}

	return userData, nil
}

func (s *userService) GetDetailUserById(id string) (res response.UserDetail, err error) {
	userData, err := s.repo.GetUserByID(id)
	if userData == nil {
		if err != nil {
			return res, errors.New(contract.ErrInternalServer)
		}
		return res, errors.New(contract.ErrRecordNotFound)
	}
	fmt.Println("USER : ", userData)
	acc, err := s.repoAccount.FindAllAccountsByUserId(userData.ID)
	if acc == nil {
		if err != nil {
			return res, errors.New(contract.ErrInternalServer)
		}
		return res, errors.New(contract.ErrRecordNotFound)
	}
	res.Email = userData.Email
	res.ID = userData.ID
	for _, accountData := range acc {
		resAccount := dtoResAccount.AccountResponse{
			ID:            accountData.ID,
			AccountNumber: accountData.AccountNumber,
			Balance:       accountData.Balance,
		}
		res.Account = append(res.Account, resAccount)
	}
	return
}
