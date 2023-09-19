package request

type TransferInquiryReq struct {
	SofAccountNumber string  `json:"sofAccountNumber" validate:"required,min=6"`
	DofAccountNumber string  `json:"dofAccountNumber" validate:"required,min=6"`
	Amount           float64 `json:"amount" validate:"required"`
}

type TransferInquiryExec struct {
	InquiryKey string `json:"inquiryKey" validate:"required,min=6"`
	Pin        string `json:"pin" validate:"required,min=6"`
}
