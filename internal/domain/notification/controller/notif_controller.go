package controller

import (
	"github.com/FaisalMashuri/my-wallet/internal/domain/notification"
	"github.com/FaisalMashuri/my-wallet/internal/domain/user"
	"github.com/FaisalMashuri/my-wallet/shared"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type notifController struct {
	service notification.NotificationService
}

func (n *notifController) GetUserNotification(ctx *fiber.Ctx) error {
	//TODO implement me
	credentialuser := ctx.Locals("credentials").(user.User)
	notifications, err := n.service.FindByUserId(credentialuser.ID)
	if err != nil {
		return err
	}
	resp := shared.SuccessResponse("success", "Success get notifications", notifications)
	return ctx.Status(http.StatusOK).JSON(resp)
}

func NewController(service notification.NotificationService) notification.NotificationController {
	return &notifController{
		service: service,
	}
}
