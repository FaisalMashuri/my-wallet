package request

type TopUpRequest struct {
	Amount        float64 `json:"amount" validate:"required"`
	UserID        string  `json:"-"`
	AccountNumber string  `json:"accountNumber" validate:"required,min=6"`
}
