package controller

import (
	"github.com/FaisalMashuri/my-wallet/internal/domain/mpin"
	"github.com/FaisalMashuri/my-wallet/internal/domain/mpin/dto/request"
	"github.com/FaisalMashuri/my-wallet/shared"
	"github.com/FaisalMashuri/my-wallet/shared/contract"
	"github.com/gofiber/fiber/v2"
)

type mPinControllerImpl struct {
	service mpin.PinService
}

func (m mPinControllerImpl) CreatePin(ctx *fiber.Ctx) error {
	//TODO implement me
	var mPinRequest request.PinRequest
	err := ctx.BodyParser(&mPinRequest)
	if err != nil {
		return fiber.NewError(400, contract.ErrBadRequest)
	}
	err = m.service.CreatePin(mPinRequest)
	if err != nil {
		return fiber.NewError(500, contract.ErrInternalServer)
	}
	resp := shared.SuccessResponse("Success", "Pin Created Successfully", nil)
	return ctx.Status(fiber.StatusCreated).JSON(resp)

}

func NewController(service mpin.PinService) mpin.PinController {
	return &mPinControllerImpl{
		service: service,
	}
}
