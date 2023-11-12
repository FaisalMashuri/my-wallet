package mpin

import (
	"github.com/FaisalMashuri/my-wallet/internal/domain/mpin/dto/request"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PinRepository interface {
	GetDB() *gorm.DB
	CreatePin(modelPin Pin) error
	FindByUserId(userId string) (Pin, error)
}

type PinService interface {
	CreatePin(modelPin request.PinRequest) error
	ValidatePin(request *request.ValidatePinReq) error
}

type PinController interface {
	CreatePin(ctx *fiber.Ctx) error
}

type Pin struct {
	gorm.Model
	ID     string `gorm:"primary_key"`
	UserID string `gorm:"user_id"`
	Pin    string `gorm:"pin"`
}

func (p *Pin) BeforeCreate(tx *gorm.DB) error {
	p.ID = uuid.NewString()
	return nil
}
