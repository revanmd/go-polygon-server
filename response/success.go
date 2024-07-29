package response

import (
	"github.com/gofiber/fiber/v2"
)

type SuccessResponse struct {
	StatusCode int         `json:"status_code"`
	StatusDesc string      `json:"status_desc"`
	Data       interface{} `json:"data"`
}

func SuccessHandler(c *fiber.Ctx, statusCode int, statusDesc string, data interface{}) error {
	response := SuccessResponse{
		StatusCode: statusCode,
		StatusDesc: statusDesc,
		Data:       data,
	}
	return c.Status(statusCode).JSON(response)
}
