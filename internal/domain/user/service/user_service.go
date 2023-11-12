package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/FaisalMashuri/my-wallet/internal/domain/account"
	dtoResAccount "github.com/FaisalMashuri/my-wallet/internal/domain/account/dto/response"
	"github.com/FaisalMashuri/my-wallet/internal/domain/mpin"
	"github.com/FaisalMashuri/my-wallet/internal/domain/mq"
	"github.com/FaisalMashuri/my-wallet/internal/domain/user"
	"github.com/FaisalMashuri/my-wallet/internal/domain/user/dto/request"
	"github.com/FaisalMashuri/my-wallet/internal/domain/user/dto/response"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"

	"github.com/FaisalMashuri/my-wallet/shared"
	"github.com/FaisalMashuri/my-wallet/shared/contract"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type userService struct {
	repo        user.UserRepository
	repoAccount account.AccountRepository
	repoPin     mpin.PinRepository
	log         *logrus.Logger
	redisClient *redis.Client
	mq          mq.MQService
}

func NewService(repo user.UserRepository, log *logrus.Logger, repoAccount account.AccountRepository, repoPin mpin.PinRepository, redis *redis.Client, mq mq.MQService) user.UserService {
	return &userService{
		repo:        repo,
		repoAccount: repoAccount,
		repoPin:     repoPin,
		log:         log,
		redisClient: redis,
		mq:          mq,
	}
}

func (s *userService) RegisterUser(userRequest *request.RegisterRequest) (*response.AuthResponse, error) {
	//ctx := context.Background()

	userModel := &user.User{
		Name:       userRequest.Name,
		Email:      userRequest.Email,
		Phone:      userRequest.Phone,
		VerifiedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}
	userData, err := s.repo.CreateUser(userModel)
	if err != nil {
		return nil, err
	}
	mPinModel := mpin.Pin{
		UserID: userData.ID,
		Pin:    userRequest.SecurityCode,
	}

	err = s.repoPin.CreatePin(mPinModel)
	if err != nil {
		return nil, err

	}
	accountModel := account.Account{
		UserID: userData.ID,
	}
	_, err = s.repoAccount.CreateAccount(accountModel)
	if err != nil {
		return nil, err
	}

	authData := user.AuthData{
		ID:           userData.ID,
		Email:        userData.Email,
		Name:         userData.Name,
		Phone:        userData.Phone,
		SecurityCode: mPinModel.Pin,
		UserType:     "regular",
		IsVerified:   true,
	}

	accessToken, err := shared.GenerateAccessToken(authData)
	if err != nil {
		return nil, err
	}

	refreshToken, err := shared.GenerateRefreshToken(authData.ID)
	if err != nil {
		return nil, err
	}
	resData := &response.AuthResponse{
		accessToken,
		refreshToken,
	}

	return resData, nil

}
func (s *userService) Login(userRequest *request.LoginRequest) (*response.AuthResponse, error) {
	userData, err := s.repo.GetUserByPhoneNumber(userRequest.Phone)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(contract.ErrRecordNotFound)
		}
		return nil, errors.New(contract.ErrUnExpected)
	}

	mPinData, err := s.repoPin.FindByUserId(userData.ID)
	if userRequest.SecurityCode != mPinData.Pin {
		return nil, errors.New(contract.ErrInvalidSecurityCode)
	}

	authData := user.AuthData{
		ID:           userData.ID,
		Email:        userData.Email,
		Name:         userData.Name,
		Phone:        userData.Phone,
		SecurityCode: mPinData.Pin,
		UserType:     "regular",
		IsVerified:   true,
	}

	accessToken, err := shared.GenerateAccessToken(authData)
	if err != nil {
		return nil, errors.New(contract.ErrGenerateToken)
	}
	refreshToken, err := shared.GenerateRefreshToken(userData.ID)
	if err != nil {
		return nil, errors.New(contract.ErrGenerateToken)
	}

	res := response.AuthResponse{
		accessToken,
		refreshToken,
	}

	return &res, nil
}

func (s *userService) GetDetailUserById(id string) (res response.UserDetail, err error) {
	var userData *user.User
	ctx := context.Background()
	val, err := s.redisClient.Get(ctx, id).Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("data tidak ditemukan di cache")
			userData, err = s.repo.GetUserByID(id)
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
			data, _ := json.Marshal(res)
			_, err = s.redisClient.Set(ctx, id, data, 15*time.Second).Result()
			if err != nil {
				s.log.Error("Error save to cache")
			}
			return res, err

		}
	}
	err = json.Unmarshal([]byte(val), &res)
	if err != nil {
		return
	}
	fmt.Println("data dari redis : ", res)

	return res, nil
}

func (s *userService) VerifyUser(verReq request.VerifiedUserRequest) error {
	ctx := context.Background()
	val, err := s.redisClient.Get(ctx, verReq.Otp).Result()
	if err != nil {
		if err == redis.Nil {
			return errors.New(contract.ErrInvalidOTP)
		}
		return errors.New(contract.ErrUnExpected)
	}
	err = s.repo.VerifyUser(val)
	if err != nil {
		return err
	}
	err = s.redisClient.Del(ctx, verReq.Otp).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) ResendOTP(userId string) error {
	ctx := context.Background()
	otp := shared.GenerateOTP()
	err := s.redisClient.Set(ctx, otp, userId, 60*time.Second).Err()
	if err != nil {
		return err
	}
	fmt.Println("OTP resend : ", otp)
	return nil
}

func (s *userService) IsPhoneNumberExist(phoneNumber string) (bool, error) {
	err := s.repo.FindPhoneNumber(phoneNumber)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, errors.New(contract.ErrUnExpected)
	}
	return true, errors.New(contract.ErrDuplicateValue)
}
