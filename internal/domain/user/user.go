package user

import (
	"github.com/FaisalMashuri/my-wallet/internal/domain/user/dto/request"
	"github.com/FaisalMashuri/my-wallet/internal/domain/user/dto/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	Login(userRequest *request.LoginRequest) (res response.LoginResponse, err error)
}

type UserController interface {
	Login(ctx *fiber.Ctx) error
	Register(ctx *fiber.Ctx) error
}

func (u *User) BeforeCreate(tx *gorm.DB) {
	u.ID = uuid.New().String()
}

func (u *User) FromLoginRequest(req request.LoginRequest) {
	u.Email = req.Email
	u.Password = req.Password
}
