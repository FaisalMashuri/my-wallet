package request

type TopUpRequest struct {
	Amount        float64 `json:"amount"`
	UserID        string  `json:"-"`
	AccountNumber string  `json:"accountNumber"`
}
