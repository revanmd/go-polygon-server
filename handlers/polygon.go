package handlers

import (
	"context"
	"polygon-server/models"
	"polygon-server/response"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PolygonHandler struct {
	DB *mongo.Database
}

func (h *PolygonHandler) GetPolygons(c *fiber.Ctx) error {
	// Parse limit and offset query parameters
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}
	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	var polygons []models.Polygon
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))

	cursor, err := h.DB.Collection("polygons").Find(context.Background(), bson.D{}, findOptions)
	if err != nil {
		return response.InternalServerError(c, err, "Error finding polygons")
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &polygons); err != nil {
		return response.InternalServerError(c, err, "Error decoding polygons")
	}

	return response.SuccessHandler(c, fiber.StatusOK, "Polygons retrieved successfully", polygons)
}
