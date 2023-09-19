package middleware

import (
	"fmt"
	"github.com/FaisalMashuri/my-wallet/shared"
	"github.com/FaisalMashuri/my-wallet/shared/contract"
	"github.com/Saucon/errcntrct"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

func NewErrorhandler(ctx *fiber.Ctx, err error) error {
	fmt.Println("ERROR HANDLING : ", err)
	if err != nil {
		dataResp, message := ErrorHandler(err)
		var errData = message.(errcntrct.ErrorData)
		resp := shared.ErrorResponse(errData.Code, "failed", fmt.Sprintf("%v", errData.Msg))
		return ctx.Status(dataResp).JSON(resp)

	}
	return nil
}

func ErrorHandler(err error) (int, interface{}) {
	var e *fiber.Error

	var code, extraError interface{}
	if errors.As(err, &e) {
		fmt.Println("ERROR : ", e)
		fmt.Print("MEssage  : ", e.Message)
		message := strings.Split(e.Message, ",")
		fmt.Println("Message split: ", message)
		fmt.Println("panjang message : ", len(message))
		if len(message) > 1 {
			code = message[0]
			extraError = message[1]
		} else {
			code = e.Message
		}

	}
	//fmt.Println("CODE : ", code)

	switch code {
	case http.StatusMethodNotAllowed:
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
	case contract.ErrUnauthorized:
		return errcntrct.ErrorMessage(http.StatusUnauthorized, "", errors.New(contract.ErrUnauthorized))
	case contract.ErrTransactionUnauthoried:
		return errcntrct.ErrorMessage(http.StatusUnauthorized, "", errors.New(contract.ErrTransactionUnauthoried))
	case contract.ErrBadRequest:
		return errcntrct.ErrorMessage(http.StatusBadRequest, "", errors.New(contract.ErrBadRequest))
	case contract.ErrInvalidPin:
		return errcntrct.ErrorMessage(http.StatusBadRequest, "", errors.New(contract.ErrInvalidPin))
	case contract.ErrLimitAccountOpen:
		return errcntrct.ErrorMessage(http.StatusInternalServerError, "", errors.New(contract.ErrLimitAccountOpen))
	case contract.ErrMandatory:
		return responseAdapter(http.StatusBadRequest, contract.ErrMandatory, extraError.(string))
	case contract.ErrFormatField:
		return responseAdapter(http.StatusBadRequest, contract.ErrFormatField, extraError.(string))
	case contract.ErrMinFormat:
		return responseAdapter(http.StatusBadRequest, contract.ErrMinFormat, extraError.(string))
	default:
		return errcntrct.ErrorMessage(9999, "", errors.New(contract.ErrUnexpectedError))
	}
}

func responseAdapter(code int, errCode string, extraError string) (int, errcntrct.ErrorData) {
	_, errData := errcntrct.ErrorMessage(code, "", errors.New(errCode))
	errData.Msg = fmt.Sprintf(errData.Msg, strings.ToLower(extraError))
	return code, errData
}
