package shared

import (
	"github.com/FaisalMashuri/my-wallet/config"
	"github.com/FaisalMashuri/my-wallet/internal/domain/user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

func GenerateAccessToken(user user.AuthData) (string, error) {
	claims := jwt.MapClaims{
		"sub":          user.ID,
		"iss":          "wallet indonesia corp",
		"iat":          time.Now().Unix(),
		"exp":          time.Now().Add(time.Hour * 32).Unix(),
		"jti":          uuid.NewString(),
		"email":        user.Email,
		"name":         user.Name,
		"phone":        user.Phone,
		"securityCode": user.SecurityCode,
		"userType":     user.UserType,
		"isVerified":   user.IsVerified,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.AppConfig.SecretKey))
	if err != nil {
		return "", err
	}
	return t, nil
}

func GenerateRefreshToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"iss": "wallet indonesia corp",
		"iat": time.Now(),
		"exp": time.Now().Add(time.Hour * 32).Unix(),
		"jti": "refresh_token_" + uuid.NewString(),
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
