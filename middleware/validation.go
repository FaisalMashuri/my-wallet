package middleware

import (
	"github.com/FaisalMashuri/my-wallet/shared/contract"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type ErrorResponse struct {
	Error       string
	FailedField string
	Tag         string
	Value       interface{}
}

func ValidateRequest(data interface{}) (string, error) {
	validate := validator.New()
	errs := validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			switch err.Tag() {
			case "email":
				return err.Field(), errors.New(contract.ErrInvalidFieldFormat)
			case "required":
				return err.Field(), errors.New(contract.ErrMandatoryField)
			case "min":
				return err.Field(), errors.New(contract.ErrMinFormat)
			}
		}
	}

	return "", nil
}
