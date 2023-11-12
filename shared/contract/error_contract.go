package contract

const (
	SuccessCode      = "0000"
	SuccessNoContent = "0200"

	ErrDuplicateValue = "0001"
	ErrRecordNotFound = "0002"
	ErrGenerateToken  = "0003"

	ErrMandatoryField          = "0100"
	ErrInvalidFieldFormat      = "0101"
	ErrCantParseBodyJSON       = "0102"
	ErrInvalidSecurityCode     = "0103"
	ErrInvalidRequestFamily    = "1000"
	ErrPasswordNotMatch        = "1001"
	ErrContextDeadlineExceeded = "1002"

	ErrEmailAlreadyRegister = "0300"
	ErrInsufficientBalance  = "0301"
	ErrLimitAccountOpen     = "0302"
	ErrInvalidPin           = "0303"
	ErrUserNotVerified      = "0304"

	ErrMethodNotAllowed        = "4050"
	ErrUrlNotFound             = "4040"
	ErrUnauthorized            = "4010"
	ErrTransactionUnauthorized = "4011"
	ErrBadRequest              = "4000"
	ErrMandatory               = "4001"
	ErrFormatField             = "4002"
	ErrMinFormat               = "4003"
	ErrInvalidTransferKey      = "4004"
	ErrInvalidOTP              = "4005"

	ErrInternalServer = "5000"

	ErrUnExpected = "9999"

	Err
	DescErrContextDeadlineExceeded = "context deadline exceeded"
)
