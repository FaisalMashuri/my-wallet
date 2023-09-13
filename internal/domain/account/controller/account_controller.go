package controller

import (
	"fmt"
	"github.com/FaisalMashuri/my-wallet/internal/domain/account"
	"github.com/FaisalMashuri/my-wallet/internal/domain/user"
	"github.com/FaisalMashuri/my-wallet/shared"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type accountControllerImpl struct {
	service account.AccountService
}

func NewController(service account.AccountService) account.AccountController {
	return &accountControllerImpl{
		service: service,
	}
}

func (c *accountControllerImpl) CreateAccount(ctx *fiber.Ctx) error {
	credential := ctx.Locals("credentials").(user.User)
	accountData, err := c.service.CreateAccount(credential.ID)
	if err != nil {
		fmt.Println("ERROR CRREATE ACCOUNT : ", err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	resp := shared.SuccessResponse("Success", "Account created", accountData)
	return ctx.Status(http.StatusCreated).JSON(resp)
}
