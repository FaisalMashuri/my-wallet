package controller

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/FaisalMashuri/my-wallet/internal/domain/notification/dto/response"
	"github.com/FaisalMashuri/my-wallet/internal/domain/sse/dto"
	"github.com/FaisalMashuri/my-wallet/internal/domain/user"
	"github.com/gofiber/fiber/v2"
)

type NotificationSseController interface {
	StreamNotifcation(ctx *fiber.Ctx) error
}

type notificationSseController struct {
	hub *dto.Hub
}

func NewNotification(hub *dto.Hub) NotificationSseController {
	return &notificationSseController{
		hub: hub,
	}
}

func (n *notificationSseController) StreamNotifcation(ctx *fiber.Ctx) error {
	ctx.Set("Content-Type", "text/event-stream")
	credentialuser := ctx.Locals("credentials").(user.User)
	fmt.Println(credentialuser)
	n.hub.NotificationChannel[credentialuser.ID] = make(chan response.NotificationDataRes)
	fmt.Println("HUB : ", n.hub.NotificationChannel[credentialuser.ID])
	ctx.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		event := fmt.Sprintf("event: %s\ndata: \n\n", "initial")
		_, _ = fmt.Fprint(w, event)
		_ = w.Flush()

		for notification := range n.hub.NotificationChannel[credentialuser.ID] {
			fmt.Println("NOTIF : ", notification)
			data, _ := json.Marshal(notification)
			fmt.Println("NOTIF : ", data)
			event := fmt.Sprintf("event: %s\ndata: %s\n\n", "notification-update", data)
			_, _ = fmt.Fprint(w, event)
			_ = w.Flush()
		}

	})
	return nil
}
