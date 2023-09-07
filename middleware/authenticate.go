package middleware

import (
	"fmt"
	"github.com/FaisalMashuri/my-wallet/internal/domain/user"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"log"
)

func NewAuthMiddleware(secret string) fiber.Handler {
	log.Println("Secret : ", secret)
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(secret)},
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(fiber.StatusUnauthorized).JSON(err)
		},
	})
}

func GetCredential(ctx *fiber.Ctx) (err error) {
	defer func() {
		if err := recover(); err != nil {
			ctx.Status(fiber.StatusInternalServerError).JSON(err)
		}
	}()
	useData := ctx.Locals("user").(*jwt.Token)
	claims := useData.Claims.(jwt.MapClaims)
	fmt.Println("CREDENTIALS: ", claims)

	credentials := user.User{
		ID:    claims["id"].(string),
		Email: claims["email"].(string),
	}
	ctx.Locals("credentials", credentials)

	return ctx.Next()
}
