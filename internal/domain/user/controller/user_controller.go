package controller

import (
	"github.com/FaisalMashuri/my-wallet/internal/domain/user"
	"github.com/FaisalMashuri/my-wallet/internal/domain/user/dto/request"
	"github.com/FaisalMashuri/my-wallet/shared"
	"github.com/FaisalMashuri/my-wallet/shared/contract"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"strconv"
)

type userController struct {
	service user.UserService
	log     *logrus.Logger
}

func NewController(service user.UserService, log *logrus.Logger) user.UserController {
	return &userController{
		service: service,
		log:     log,
	}
}

func (c *userController) Login(ctx *fiber.Ctx) error {
	var loginReq request.LoginRequest
	err := ctx.BodyParser(&loginReq)
	if err != nil {
		c.log.WithField("error", err.Error()).Fatal("error parsing request body")
		errCode, _ := strconv.Atoi(err.Error())
		return fiber.NewError(errCode, contract.ErrInvalidRequestFamily)
	}
	res, err := c.service.Login(&loginReq)
	if err != nil {
		c.log.WithField("error", err.Error()).Info("login failed " + err.Error())
		errCode, _ := strconv.Atoi(err.Error())
		return fiber.NewError(errCode, err.Error())
	}
	resp := shared.SuccessResponse("Success", "login success", res)
	return ctx.Status(200).JSON(resp)
}

func (c *userController) Register(ctx *fiber.Ctx) error {
	var regisRequest request.RegisterRequest
	err := ctx.BodyParser(&regisRequest)
	if err != nil {
		c.log.WithField("error", err.Error()).Fatal("error parsing request body")
		errCode, _ := strconv.Atoi(err.Error())
		return fiber.NewError(errCode, contract.ErrInvalidRequestFamily)
	}
	userData, err := c.service.RegisterUser(&regisRequest)
	if err != nil {
		c.log.WithField("error", err.Error()).Info("Registration failed " + err.Error())
		errCode, _ := strconv.Atoi(err.Error())
		return fiber.NewError(errCode, err.Error())
	}
	resp := shared.SuccessResponse("Succes", "Succes registration user", userData.ToRegisterResponse())
	return ctx.Status(200).JSON(resp)
}
