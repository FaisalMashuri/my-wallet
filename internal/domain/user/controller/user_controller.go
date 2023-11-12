package controller

import (
	"fmt"
	"github.com/FaisalMashuri/my-wallet/internal/domain/user"
	"github.com/FaisalMashuri/my-wallet/internal/domain/user/dto/request"
	"github.com/FaisalMashuri/my-wallet/middleware"
	"github.com/FaisalMashuri/my-wallet/shared"
	"github.com/FaisalMashuri/my-wallet/shared/constant"
	"github.com/FaisalMashuri/my-wallet/shared/contract"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
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
	fieldErr, err := middleware.ValidateRequest(loginReq)
	fmt.Println("FIELD ERR  : ", fieldErr)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, fmt.Sprintf("%s,%s", err.Error(), fieldErr))
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
	fieldErr, err := middleware.ValidateRequest(regisRequest)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, fmt.Sprintf("%s,%s", err.Error(), fieldErr))
	}
	res, err := c.service.RegisterUser(&regisRequest)
	if err != nil {
		c.log.WithField("error", err.Error()).Info("Registration failed " + err.Error())
		return errors.New(err.Error())
	}
	resp := shared.SuccessResponse("Succes", "Succes registration user", res)
	return middleware.ResponseSuccess(ctx, contract.SuccessCode, resp)
}

func (c *userController) GetDetailUserJWT(ctx *fiber.Ctx) error {
	credentialuser := ctx.Locals("credentials").(user.User)
	fmt.Println("CREDENTIAL : ", credentialuser)
	res, err := c.service.GetDetailUserById(credentialuser.ID)
	if err != nil {
		return fiber.NewError(500, err.Error())
	}
	resp := shared.SuccessResponse("Success", "Success get detail user", res)
	return ctx.Status(200).JSON(resp)
}

func (c *userController) VerifyUser(ctx *fiber.Ctx) error {
	fmt.Println("VERIFYING USER")
	var verReq request.VerifiedUserRequest
	err := ctx.BodyParser(&verReq)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, contract.ErrBadRequest)
	}
	err = c.service.VerifyUser(verReq)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}
	resp := shared.SuccessResponse("Success", "User has been verified", nil)
	return ctx.Status(http.StatusOK).JSON(resp)
}

func (c *userController) ResendOTP(ctx *fiber.Ctx) error {
	fmt.Println("ResendOTP")
	userId := ctx.Query("user")
	if userId == constant.EmptyString {
		return fiber.NewError(http.StatusBadRequest, contract.ErrBadRequest)
	}
	err := c.service.ResendOTP(userId)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, contract.ErrBadRequest)
	}
	resp := shared.SuccessResponse("Success", "OTP has been resent", nil)
	return ctx.Status(http.StatusOK).JSON(resp)
}

func (c *userController) CheckPhoneNumberExist(ctx *fiber.Ctx) error {
	var req request.CheckPhoneNumberRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		return errors.New(contract.ErrCantParseBodyJSON)
	}
	fmt.Println("REQUEST : ", req.Phone)
	isExist, err := c.service.IsPhoneNumberExist(req.Phone)
	if err != nil || isExist {
		return err
	}

	return middleware.ResponseSuccess(ctx, contract.SuccessCode, "Phone number available")

}
