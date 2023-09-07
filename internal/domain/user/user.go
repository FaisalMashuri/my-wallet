package user

import (
	"github.com/FaisalMashuri/my-wallet/internal/domain/user/dto/request"
	"github.com/FaisalMashuri/my-wallet/internal/domain/user/dto/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       string `gorm:"primary_key"`
	Email    string
	Password string
}

type UserRepository interface {
	FindUserByEmail(email string) (*User, error)
	CreateUser(user *User) (*User, error)
	UpdateUser(updatedUser *User) (*User, error)
	GetUserByID(id string) (*User, error)
	GetAllUser() ([]*User, error)
	DeleteUser(id string) error
}

type UserService interface {
	Login(userRequest *request.LoginRequest) (res *response.LoginResponse, err error)
	RegisterUser(userRequest *request.RegisterRequest) (user *User, err error)
}

type UserController interface {
	Login(ctx *fiber.Ctx) error
	Register(ctx *fiber.Ctx) error
}

func (u *User) FromRegistRequest(req *request.RegisterRequest) User {
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
	return User{
		Email:    req.Email,
		Password: string(hashedPass),
	}
}

func (u *User) ToRegisterResponse() response.RegisterResponse {
	return response.RegisterResponse{
		ID:        u.ID,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()

	return
}
