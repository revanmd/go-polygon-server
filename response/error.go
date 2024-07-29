package response

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandler is a centralized function for handling errors
func ErrorHandler(c *fiber.Ctx, err error, status int, message string) error {
	if err != nil {
		log.Println(message, err)
	} else {
		log.Println(message)
	}
	return c.Status(status).JSON(fiber.Map{"error": message})
}

// BadRequest handles 400 errors
func BadRequest(c *fiber.Ctx, err error, message string) error {
	return ErrorHandler(c, err, fiber.StatusBadRequest, message)
}

// ValidationError handles validation errors
func ValidationError(c *fiber.Ctx, message string) error {
	return ErrorHandler(c, nil, fiber.StatusBadRequest, message)
}

// InternalServerError handles 500 errors
func InternalServerError(c *fiber.Ctx, err error, message string) error {
	return ErrorHandler(c, err, fiber.StatusInternalServerError, message)
}
