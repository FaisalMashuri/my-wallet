package shared

import (
	"github.com/FaisalMashuri/my-wallet/config"
	"github.com/FaisalMashuri/my-wallet/internal/domain/user"
	"github.com/golang-jwt/jwt/v5"
	"math/rand"
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

		return "", err
	}
	return t, nil
}

func GenerateAccountNumber() string {
	const charset = "0123456789"
	const codeLength = 6
	rand.Seed(time.Now().UnixNano())

	// Menghasilkan karakter acak
	randomPart := make([]byte, codeLength)
	for i := 0; i < len(randomPart); i++ {
		randomPart[i] = charset[rand.Intn(len(charset))]
	}

	return string(randomPart)
}

func GenerateInquiryKey() string {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const codeLength = 8
	rand.Seed(time.Now().UnixNano())

	// Menghasilkan karakter acak
	randomPart := make([]byte, codeLength)
	for i := 0; i < len(randomPart); i++ {
		randomPart[i] = charset[rand.Intn(len(charset))]
	}

	return string(randomPart)
}

func GenerateOTP() string {
	const charset = "0123456789"
	const codeLength = 6

	rand.New(rand.NewSource(time.Now().UnixNano()))
	randomPart := make([]byte, codeLength)
	for i := 0; i < len(randomPart); i++ {
		randomPart[i] = charset[rand.Intn(len(charset))]
	}
	return string(randomPart)
}
