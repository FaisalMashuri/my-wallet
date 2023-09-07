package repository

import (
	"github.com/FaisalMashuri/my-wallet/internal/domain/notification"
	"gorm.io/gorm"
)

type notificationRepository struct {
	db *gorm.DB
}

func (n *notificationRepository) FindByUserId(userId string) (notif []notification.Notification, err error) {
	//TODO implement me
	err = n.db.Debug().Find(&notif).Limit(15).Order("created_at").Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return []notification.Notification{}, nil
		}
		return []notification.Notification{}, err
	}
	return notif, nil
}

func (n *notificationRepository) InsertNotification(notification *notification.Notification) error {
	//TODO implement me
	err := n.db.Debug().Create(&notification).Error
	if err != nil {
		return err
	}
	return nil
}

func (n notificationRepository) UpdateNotification(notification *notification.Notification) error {
	//TODO implement me
	err := n.db.Debug().Save(&notification).Error
	if err != nil {
		return err
	}
	return nil

}

func NewRepository(db *gorm.DB) notification.NotificationRepository {
	return &notificationRepository{
		db: db,
	}
}
