package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/FaisalMashuri/my-wallet/internal/domain/account"
	dtoResAccount "github.com/FaisalMashuri/my-wallet/internal/domain/account/dto/response"
	"github.com/FaisalMashuri/my-wallet/internal/domain/mq"
	"github.com/FaisalMashuri/my-wallet/internal/domain/user"
	"github.com/FaisalMashuri/my-wallet/internal/domain/user/dto/request"
	"github.com/FaisalMashuri/my-wallet/internal/domain/user/dto/response"
	"github.com/redis/go-redis/v9"
	"time"

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
	redisClient *redis.Client
	mq          mq.MQService
}

func NewService(repo user.UserRepository, log *logrus.Logger, repoAccount account.AccountRepository, redis *redis.Client, mq mq.MQService) user.UserService {
	return &userService{
		repo:        repo,
		repoAccount: repoAccount,
		log:         log,
		redisClient: redis,
		mq:          mq,
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
	if !userData.VerifiedAt.Valid {
		return nil, errors.New(contract.ErrUserNotVerified)
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
	ctx := context.Background()
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
	otp := shared.GenerateOTP()
	err = s.redisClient.Set(ctx, otp, userModel.ID, 60*time.Second).Err()
	if err != nil {
		return nil, err
	}
	fmt.Println("OTP : ", otp)
	payloadBody := map[string]interface{}{
		"email": userModel.Email,
		"otp":   otp,
	}

	pBody := mq.MessageQueueEmitter{
		Title:   "Register user",
		Payload: payloadBody,
	}
	payload, err := json.Marshal(pBody)

	err = s.mq.SendData("user.register", payload)

	return userData, nil
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
		return errors.New(contract.ErrUnexpectedError)
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
