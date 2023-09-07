package middleware

import (
	"fmt"
	"github.com/FaisalMashuri/my-wallet/shared"
	"github.com/FaisalMashuri/my-wallet/shared/contract"
	"github.com/Saucon/errcntrct"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"net/http"
)

func NewErrorhandler(ctx *fiber.Ctx, err error) error {
	if err != nil {
		dataResp, message := ErrorHandler(err)
		var errData = message.(errcntrct.ErrorData)
		resp := shared.ErrorResponse(errData.Code, "failed", fmt.Sprintf("%v", errData.Msg))
		return ctx.Status(dataResp).JSON(resp)

	}
	return nil
}

func ErrorHandler(error error) (int, interface{}) {
	var e *fiber.Error
	var code interface{}

	if errors.As(error, &e) {
		code = e.Message

	}
	fmt.Println("CODE : ", code)

	switch code {
	case http.StatusMethodNotAllowed:
		fmt.Println("hai")
		return errcntrct.ErrorMessage(http.StatusMethodNotAllowed, "", errors.New(contract.ErrMethodNotAllowed))
	case http.StatusNotFound:
		return errcntrct.ErrorMessage(http.StatusNotFound, "", errors.New(contract.ErrUrlNotFound))
	case contract.ErrRecordNotFound:
		return errcntrct.ErrorMessage(http.StatusInternalServerError, "", errors.New(contract.ErrRecordNotFound))
	case contract.ErrEmailAlreadyRegister:
		return errcntrct.ErrorMessage(http.StatusInternalServerError, "", errors.New(contract.ErrEmailAlreadyRegister))
	case contract.ErrInvalidRequestFamily:
		return errcntrct.ErrorMessage(http.StatusBadRequest, "", errors.New(contract.ErrInvalidRequestFamily))
	case contract.ErrPasswordNotMatch:
		return errcntrct.ErrorMessage(http.StatusUnauthorized, "", errors.New(contract.ErrPasswordNotMatch))
	case contract.ErrInternalServer:
		return errcntrct.ErrorMessage(http.StatusInternalServerError, "", errors.New(contract.ErrInternalServer))
	case contract.ErrContextDeadlineExceeded:
		return errcntrct.ErrorMessage(http.StatusGatewayTimeout, "", errors.New(contract.ErrContextDeadlineExceeded))
	}
	if code == nil {
		code = "9999"
	}

	return errcntrct.ErrorMessage(9999, "", errors.New(contract.ErrUnexpectedError))

}
