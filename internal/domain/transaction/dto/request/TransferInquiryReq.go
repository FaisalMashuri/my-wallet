package request

type TransferInquiryReq struct {
	SofAccountNumber string  `json:"sofAccountNumber"`
	DofAccountNumber string  `json:"dofAccountNumber"`
	Amount           float64 `json:"amount"`
}

type TransferInquiryExec struct {
	InquiryKey string `json:"inquiryKey"`
}
