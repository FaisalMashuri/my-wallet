package controller

import (
	"github.com/FaisalMashuri/my-wallet/internal/domain/user"
	"github.com/gofiber/fiber/v2"
)

type userController struct {
	service user.UserService
}

func NewController(service user.UserService) user.UserController {
	return &userController{
		service: service,
	}
}

func (c *userController) Login(ctx *fiber.Ctx) error {
	return ctx.Status(200).JSON("OKe")
}

func (c *userController) Register(ctx *fiber.Ctx) error {
	return ctx.Status(200).JSON("OKe")
}
