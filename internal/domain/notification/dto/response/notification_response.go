package response

import "time"

type NotificationDataRes struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	IsRead    int8      `json:"isRead"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}
