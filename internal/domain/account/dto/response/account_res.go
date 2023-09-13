package response

type AccountResponse struct {
	ID            string  `json:"account_id"`
	AccountNumber string  `json:"account_number"`
	Balance       float64 `json:"balance"`
}
