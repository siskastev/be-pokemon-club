package middleware

import (
	"be-pokemon-club/internal/helpers/jwt"
	"be-pokemon-club/internal/models"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	SecretKey string
}

func NewAuthMiddleware(secretKey string) *AuthMiddleware {
	return &AuthMiddleware{SecretKey: secretKey}
}

func (a *AuthMiddleware) AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
		}

		authToken := strings.Split(authHeader, " ")[1]
		user, err := jwt.VerifyAndExtractUserFromJWT(authToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
		}

		c.Locals("user", user)

		return c.Next()
	}
}

func (a *AuthMiddleware) IsBoss() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, found := c.Locals("user").(*models.UserResponse)
		if !found {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "halloo"})
		}

		if user.Role != "bos" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized Role Access"})
		}

		return c.Next()
	}
}
