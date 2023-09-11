package topup

import (
	"github.com/FaisalMashuri/my-wallet/internal/domain/topup/dto/request"
	"github.com/FaisalMashuri/my-wallet/internal/domain/topup/dto/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TopUpRepository interface {
	FindById(id string) (*TopUp, error)
	Insert(t *TopUp) error
	Update(t *TopUp) error
}

type TopUpService interface {
	ConfirmedTopUp(id string) error
	InitializeTopUp(req request.TopUpRequest) (response.TopUpResponnse, error)
}

type TopUpController interface {
	InitializeTopUp(ctx *fiber.Ctx) error
}

type TopUp struct {
	gorm.Model
	ID      string  `gorm:"primary_key"`
	UserID  string  `gorm:"user_id"`
	Status  int     `gorm:"status"`
	Amount  float64 `gorm:"amount"`
	SnapURL string  `gorm:"snap_url"`
}

func (t *TopUp) BeforeCreate() error {
	t.ID = uuid.New().String()
	return nil
}
