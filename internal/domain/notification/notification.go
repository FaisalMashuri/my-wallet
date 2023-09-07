package notification

import (
	"github.com/FaisalMashuri/my-wallet/internal/domain/notification/dto/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	FindByUserId(userId string) ([]Notification, error)
	InsertNotification(notification *Notification) error
	UpdateNotification(notification *Notification) error
}

type NotificationService interface {
	FindByUserId(userId string) ([]response.NotificationDataRes, error)
}

type NotificationController interface {
	GetUserNotification(ctx *fiber.Ctx) error
}

type Notification struct {
	gorm.Model
	ID     string `gorm:"primary_key"`
	UserID string `gorm:"user_id"`
	Title  string `gorm:"title"`
	Body   string `gorm:"body"`
	IsRead int8   `gorm:"is_read"`
	Status string `gorm:"status"`
}

func (n *Notification) BeforeCreate(db *gorm.DB) error {
	n.ID = uuid.New().String()
	return nil
}
