package response

import (
	responseAccoount "github.com/FaisalMashuri/my-wallet/internal/domain/account/dto/response"
	"time"
)

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RegisterResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserDetail struct {
	ID      string                             `json:"id"`
	Email   string                             `json:"email"`
	Account []responseAccoount.AccountResponse `json:"account"`
}
