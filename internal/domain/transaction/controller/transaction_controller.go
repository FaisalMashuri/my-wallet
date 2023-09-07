package controller

import (
	"github.com/FaisalMashuri/my-wallet/internal/domain/transaction"
	"github.com/FaisalMashuri/my-wallet/internal/domain/transaction/dto/request"
	"github.com/FaisalMashuri/my-wallet/shared"
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
		return err
	}
	res, err := t.service.TranferInquiry(inquiryReq, ctx)
	if err != nil {
		return err
	}
	resp := shared.SuccessResponse("Success", "Transfer Inquiry berhaisl", res)
	return ctx.Status(http.StatusOK).JSON(resp)
	panic("implement me")
}

func (t *transactionController) TransferExec(ctx *fiber.Ctx) error {
	//TODO implement me
	var inquiryExecReq request.TransferInquiryExec
	err := ctx.BodyParser(&inquiryExecReq)
	if err != nil {
		return err
	}
	err = t.service.TransferInquiryExec(inquiryExecReq, ctx)
	if err != nil {
		return err
	}
	return ctx.JSON("OK")
	panic("implement me")
}

func NewController(service transaction.TransactionService) transaction.TransactionController {
	return &transactionController{
		service: service,
	}
}
