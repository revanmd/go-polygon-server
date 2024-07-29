package middleware

import (
	"polygon-server/utils"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

func JWTMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   utils.GetJWTSecret(),
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Unauthorized",
	})
}
