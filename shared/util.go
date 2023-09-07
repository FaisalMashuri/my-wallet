package shared

import (
	"github.com/FaisalMashuri/my-wallet/config"
	"github.com/FaisalMashuri/my-wallet/infrastructure"
	"github.com/FaisalMashuri/my-wallet/internal/domain/user"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GenerateAccessToken(user *user.User) (string, error) {
	claims := jwt.MapClaims{
		"id":        user.ID,
		"email":     user.Email,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
		"exp":       time.Now().Add(time.Hour * 1).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.AppConfig.SecretKey))
	if err != nil {
		infrastructure.Log.Error("Error signing jwt")
		return "", err
	}
	return t, nil
}

func GenerateRefreshToken(user *user.User) (string, error) {
	claims := jwt.MapClaims{
		"id":        user.ID,
		"email":     user.Email,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
		"exp":       time.Now().Add(time.Hour * 32).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.AppConfig.SecretKey))
	if err != nil {
		infrastructure.Log.Error("Error signing jwt")
		return "", err
	}
	return t, nil
}
