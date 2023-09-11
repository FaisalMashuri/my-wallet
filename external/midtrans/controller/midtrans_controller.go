package controller

import (
	"fmt"
	midtrans_ext "github.com/FaisalMashuri/my-wallet/external/midtrans"
	"github.com/FaisalMashuri/my-wallet/internal/domain/topup"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
)

type midTransController struct {
	service      midtrans_ext.MidtransService
	topUpService topup.TopUpService
}

func (m midTransController) PaymentHandlerNotification(ctx *fiber.Ctx) error {
	//TODO implement me
	var notificationPayload map[string]interface{}
	if err := ctx.BodyParser(&notificationPayload); err != nil {
		return err
	}
	fmt.Println("Notif payload : ", notificationPayload)
	orderId, exists := notificationPayload["order_id"].(string)
	if !exists {
		return ctx.Status(http.StatusBadRequest).JSON("not exists order_id")
	}

	success, err := m.service.VerifyPayment(orderId)
	log.Println("Success verify payment: ", success)
	if err != nil {
		log.Println("error : ", err)
		return err
	}
	log.Println("ORDER ID :", orderId)
	if success {
		_ = m.topUpService.ConfirmedTopUp(orderId)
		return ctx.Status(http.StatusOK).JSON("OK")
	}
	return ctx.Status(http.StatusBadRequest).JSON("ADA ERROR")
}

func NewController(service midtrans_ext.MidtransService, topUpService topup.TopUpService) midtrans_ext.MidtransController {
	return &midTransController{
		service:      service,
		topUpService: topUpService,
	}
}
