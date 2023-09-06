package service

import (
	"github.com/FaisalMashuri/my-wallet/internal/domain/user"
	"github.com/FaisalMashuri/my-wallet/internal/domain/user/dto/request"
	"github.com/FaisalMashuri/my-wallet/internal/domain/user/dto/response"
)

type userService struct {
	repo user.UserRepository
}

func NewService(repo user.UserRepository) user.UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) Login(userRequest *request.LoginRequest) (res response.LoginResponse, err error) {
	return res, nil
}
