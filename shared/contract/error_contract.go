package contract

const (
	ErrInvalidRequestFamily    = "1000"
	ErrPasswordNotMatch        = "1001"
	ErrContextDeadlineExceeded = "1002"

	ErrRecordNotFound       = "0299"
	ErrEmailAlreadyRegister = "0300"
	ErrInsuficentBalance    = "0301"
	ErrLimitAccountOpen     = "0302"
	ErrInvalidPin           = "0303"

	ErrMethodNotAllowed       = "4050"
	ErrUrlNotFound            = "4040"
	ErrUnauthorized           = "4010"
	ErrTransactionUnauthoried = "4011"
	ErrBadRequest             = "4000"

	ErrInternalServer = "5000"

	ErrUnexpectedError = "9999"

	Err
	DescErrContextDeadlineExceeded = "context deadline exceeded"
)
