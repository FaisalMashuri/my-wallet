package request

type PinRequest struct {
	UserId string `json:"user_id"`
	Pin    string `json:"pin"`
}
