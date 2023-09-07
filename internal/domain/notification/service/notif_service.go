package service

import (
	"github.com/FaisalMashuri/my-wallet/internal/domain/notification"
	"github.com/FaisalMashuri/my-wallet/internal/domain/notification/dto/response"
)

type notificationService struct {
	repo notification.NotificationRepository
}

func (n notificationService) FindByUserId(userId string) ([]response.NotificationDataRes, error) {
	notifications, err := n.repo.FindByUserId(userId)
	if err != nil {
		return nil, err
	}
	var result []response.NotificationDataRes
	for _, notificationData := range notifications {
		result = append(result, response.NotificationDataRes{
			ID:        notificationData.ID,
			Title:     notificationData.Title,
			Body:      notificationData.Body,
			Status:    notificationData.Status,
			IsRead:    notificationData.IsRead,
			CreatedAt: notificationData.CreatedAt,
		})
	}

	if result == nil {
		result = make([]response.NotificationDataRes, 0)
	}

	return result, nil
}

func NewService(repo notification.NotificationRepository) notification.NotificationService {
	return &notificationService{
		repo: repo,
	}
}
