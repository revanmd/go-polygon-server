package response

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func ErrorHandler(c *fiber.Ctx, err error, status int, message string) error {
	if err != nil {
		logrus.Error(message, err)
	} else {
		logrus.Error(message)
	}
	return c.Status(status).JSON(fiber.Map{"error": message})
}

func BadRequest(c *fiber.Ctx, err error, message string) error {
	return ErrorHandler(c, err, fiber.StatusBadRequest, message)
}

func ValidationError(c *fiber.Ctx, message string) error {
	return ErrorHandler(c, nil, fiber.StatusBadRequest, message)
}

func InternalServerError(c *fiber.Ctx, err error, message string) error {
	return ErrorHandler(c, err, fiber.StatusInternalServerError, message)
}
