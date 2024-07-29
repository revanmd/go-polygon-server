package middleware

import (
	"polygon-server/response"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func CheckMongoDBConnection(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if db == nil {
			return response.InternalServerError(c, nil, "MongoDB not initialized")
		}
		return c.Next()
	}
}
