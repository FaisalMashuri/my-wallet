package midtrans_ext

import (
	"github.com/FaisalMashuri/my-wallet/internal/domain/topup"
	"github.com/gofiber/fiber/v2"
)

type MidtransService interface {
	GenerateSnapURL(t *topup.TopUp) error
	VerifyPayment(orderId string) (bool, error)
}

type MidtransController interface {
	PaymentHandlerNotification(ctx *fiber.Ctx) error
}
