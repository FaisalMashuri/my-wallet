package controller

import (
	"github.com/FaisalMashuri/my-wallet/internal/domain/transaction"
	"github.com/FaisalMashuri/my-wallet/internal/domain/transaction/dto/request"
	"github.com/FaisalMashuri/my-wallet/shared"
	"github.com/FaisalMashuri/my-wallet/shared/contract"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type transactionController struct {
	service transaction.TransactionService
}

func (t *transactionController) TransferInquiry(ctx *fiber.Ctx) error {
	//TODO implement me
	var inquiryReq request.TransferInquiryReq
	err := ctx.BodyParser(&inquiryReq)
	if err != nil {
		return fiber.NewError(400, contract.ErrBadRequest)

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
	var inquiryExecReq request.TransferInquiryExec
	err := ctx.BodyParser(&inquiryExecReq)
	if err != nil {
		return fiber.NewError(400, contract.ErrBadRequest)
	}
	err = t.service.TransferInquiryExec(inquiryExecReq, ctx)
	if err != nil {
		if err.Error() == contract.ErrRecordNotFound {
			return fiber.NewError(404, err.Error())
		}
		return fiber.NewError(500, err.Error())
	}
	return ctx.JSON("OK")

}

func NewController(service transaction.TransactionService) transaction.TransactionController {
	return &transactionController{
		service: service,
	}
}
