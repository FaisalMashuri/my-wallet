package controller

import (
	"fmt"
	"github.com/FaisalMashuri/my-wallet/internal/domain/mpin"
	mPinRequest "github.com/FaisalMashuri/my-wallet/internal/domain/mpin/dto/request"
	"github.com/FaisalMashuri/my-wallet/internal/domain/transaction"
	"github.com/FaisalMashuri/my-wallet/internal/domain/transaction/dto/request"
	"github.com/FaisalMashuri/my-wallet/internal/domain/user"
	"github.com/FaisalMashuri/my-wallet/middleware"
	"github.com/FaisalMashuri/my-wallet/shared"
	"github.com/FaisalMashuri/my-wallet/shared/contract"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type transactionController struct {
	service     transaction.TransactionService
	servicemPin mpin.PinService
}

func (t *transactionController) TransferInquiry(ctx *fiber.Ctx) error {
	//TODO implement me
	var inquiryReq request.TransferInquiryReq
	err := ctx.BodyParser(&inquiryReq)
	if err != nil {
		return fiber.NewError(400, contract.ErrBadRequest)

	}
	fieldErr, err := middleware.ValidateRequest(inquiryReq)
	fmt.Println("FIELD Error : ", fieldErr)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, fmt.Sprintf("%s,%s", err.Error(), fieldErr))
	}
	res, err := t.service.TranferInquiry(inquiryReq, ctx)
	if err != nil {
		if err.Error() == contract.ErrRecordNotFound {
			return fiber.NewError(404, err.Error())
		}
		return fiber.NewError(500, err.Error())
	}
	resp := shared.SuccessResponse("Success", "Transfer Inquiry berhaisl", res)
	return ctx.Status(http.StatusOK).JSON(resp)

}

func (t *transactionController) TransferExec(ctx *fiber.Ctx) error {
	//TODO implement me
	var mpinReq mPinRequest.ValidatePinReq
	credentialuser := ctx.Locals("credentials").(user.User)

	var inquiryExecReq request.TransferInquiryExec
	err := ctx.BodyParser(&inquiryExecReq)
	if err != nil {
		return fiber.NewError(400, contract.ErrBadRequest)
	}
	fieldErr, err := middleware.ValidateRequest(inquiryExecReq)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, fmt.Sprintf("%s,%s", err.Error(), fieldErr))
	}

	fmt.Println("INQUIRY KEY : ", inquiryExecReq)
	mpinReq.Pin = inquiryExecReq.Pin
	mpinReq.UserId = credentialuser.ID
	err = t.servicemPin.ValidatePin(&mpinReq)
	if err != nil {
		return fiber.NewError(400, contract.ErrInvalidPin)
	}
	err = t.service.TransferInquiryExec(inquiryExecReq, ctx)

	if err != nil {
		if err.Error() == contract.ErrRecordNotFound {
			return fiber.NewError(404, err.Error())
		}
		return fiber.NewError(500, err.Error())
	}
	resp := shared.SuccessResponse("Success", "Transfer successfully", nil)
	return ctx.Status(http.StatusOK).JSON(resp)

}

func NewController(service transaction.TransactionService, servicePin mpin.PinService) transaction.TransactionController {
	return &transactionController{
		service:     service,
		servicemPin: servicePin,
	}
}
