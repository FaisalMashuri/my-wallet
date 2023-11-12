package response

import (
	responseAccoount "github.com/FaisalMashuri/my-wallet/internal/domain/account/dto/response"
)

type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type UserDetail struct {
	ID      string                             `json:"id"`
	Email   string                             `json:"email"`
	Account []responseAccoount.AccountResponse `json:"account"`
}
