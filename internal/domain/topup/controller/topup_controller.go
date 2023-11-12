package controller

import (
	"github.com/FaisalMashuri/my-wallet/internal/domain/topup"
	"github.com/FaisalMashuri/my-wallet/internal/domain/topup/dto/request"
	"github.com/FaisalMashuri/my-wallet/internal/domain/user"
	"github.com/FaisalMashuri/my-wallet/shared"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type topUpController struct {
	topUpService topup.TopUpService
}

func (t topUpController) InitializeTopUp(ctx *fiber.Ctx) error {
	//TODO implement me
	var req request.TopUpRequest
	credentialuser := ctx.Locals("credentials").(user.User)

	if err := ctx.BodyParser(&req); err != nil {
		return err
	}
	//fieldErr, err := middleware.ValidateRequest(req)
	//if err != nil {
	//	return fiber.NewError(http.StatusBadRequest, fmt.Sprintf("%s,%s", err.Error(), fieldErr))
	//}
	req.UserID = credentialuser.ID
	res, err := t.topUpService.InitializeTopUp(req)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	resp := shared.SuccessResponse("Success", "topup berhasil", res)
	return ctx.Status(http.StatusOK).JSON(resp)
}

func NewController(topUpService topup.TopUpService) topup.TopUpController {
	return &topUpController{
		topUpService: topUpService,
	}
}
