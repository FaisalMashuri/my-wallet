package user

import (
	"database/sql"
	"github.com/FaisalMashuri/my-wallet/internal/domain/user/dto/request"
	"github.com/FaisalMashuri/my-wallet/internal/domain/user/dto/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID         string `gorm:"primary_key"`
	Email      string `gorm:"unique"`
	Name       string
	Password   string
	Phone      string       `gorm:"uniqueIndex"`
	VerifiedAt sql.NullTime `gorm:"default:null"`
}

// For mapping jwt token
type AuthData struct {
	ID           string
	Email        string
	Phone        string
	Name         string
	SecurityCode string
	TotalBalance float64
	UserType     string
	IsVerified   bool
}

type UserRepository interface {
	GetDB() *gorm.DB
	FindUserByEmail(email string) (*User, error)
	CreateUser(user *User) (*User, error)
	UpdateUser(updatedUser *User) (*User, error)
	GetUserByID(id string) (*User, error)
	GetAllUser() ([]*User, error)
	DeleteUser(id string) error
	VerifyUser(id string) error
	FindPhoneNumber(phoneNumber string) error
	GetUserByPhoneNumber(phoneNumber string) (*User, error)
}

type UserService interface {
	Login(userRequest *request.LoginRequest) (res *response.AuthResponse, err error)
	RegisterUser(userRequest *request.RegisterRequest) (res *response.AuthResponse, err error)
	GetDetailUserById(id string) (user response.UserDetail, err error)
	VerifyUser(verReq request.VerifiedUserRequest) error
	ResendOTP(userId string) error
	IsPhoneNumberExist(phoneNumber string) (bool, error)
}

type UserController interface {
	Login(ctx *fiber.Ctx) error
	Register(ctx *fiber.Ctx) error
	GetDetailUserJWT(ctx *fiber.Ctx) error
	VerifyUser(ctx *fiber.Ctx) error
	ResendOTP(ctx *fiber.Ctx) error
	CheckPhoneNumberExist(ctx *fiber.Ctx) error
}

func (u *User) FromRegistRequest(req *request.RegisterRequest) User {
	//hashedPass, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
	return User{
		Email: req.Email,
		//Password: string(hashedPass),
		Phone: req.Phone,
	}
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()

	return
}
