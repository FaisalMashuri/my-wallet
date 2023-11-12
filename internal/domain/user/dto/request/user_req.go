package request

type LoginRequest struct {
	Phone        string `json:"phone" validate:"required"`
	SecurityCode string `json:"securityCode" validate:"required"`
}

type CheckPhoneNumberRequest struct {
	Phone string `json:"phone" validate:"required"`
}

type RegisterRequest struct {
	Name         string `json:"name" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	Phone        string `json:"phone" validate:"required"`
	SecurityCode string `json:"securityCode" validate:"required"`
}

type VerifiedUserRequest struct {
	Otp string `json:"otp"`
}
