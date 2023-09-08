package dto

import (
	notifData "github.com/FaisalMashuri/my-wallet/internal/domain/notification/dto/response"
)

type Hub struct {
	NotificationChannel map[string]chan notifData.NotificationDataRes
}
