package handlers

import (
	"context"
	"strconv"

	"polygon-server/models"
	"polygon-server/response"

	"github.com/gofiber/fiber/v2"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/simplify"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PolygonHandler struct {
	DB *mongo.Database
}

func (h *PolygonHandler) GetPolygons(c *fiber.Ctx) error {
	lat, err := strconv.ParseFloat(c.Query("lat"), 64)
	if err != nil {
		return response.BadRequest(c, err, "Invalid latitude")
	}

	long, err := strconv.ParseFloat(c.Query("long"), 64)
	if err != nil {
		return response.BadRequest(c, err, "Invalid longitude")
	}

	radius, err := strconv.ParseFloat(c.Query("radius"), 64)
	if err != nil || radius <= 0 {
		return response.BadRequest(c, err, "Invalid radius")
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}
	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	simplifyTolerance, err := strconv.ParseFloat(c.Query("simplify", "0"), 64)
	if err != nil || simplifyTolerance < 0 {
		simplifyTolerance = 0
	}

	var polygons []models.Polygon
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))

	// Geospatial query
	filter := bson.M{
		"geometry": bson.M{
			"$geoWithin": bson.M{
				"$centerSphere": []interface{}{
					[]float64{long, lat},
					radius / 6378100.0, // Radius in radians (6378100 is the approximate radius of Earth in meters)
				},
			},
		},
	}

	cursor, err := h.DB.Collection("polygons").Find(context.Background(), filter, findOptions)
	if err != nil {
		return response.InternalServerError(c, err, "Error finding polygons")
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &polygons); err != nil {
		return response.InternalServerError(c, err, "Error decoding polygons")
	}

	// Simplify the polygons if requested
	if simplifyTolerance > 0 {
		for i, polygon := range polygons {
			simplifiedCoordinates := simplifyPolygon(polygon.Geometry.Coordinates, simplifyTolerance)
			polygons[i].Geometry.Coordinates = simplifiedCoordinates
		}
	}

	return response.SuccessHandler(c, fiber.StatusOK, "Polygons retrieved successfully", polygons)
}

// simplifyPolygon simplifies the coordinates of a polygon using the given tolerance
func simplifyPolygon(coordinates [][][][]float64, tolerance float64) [][][][]float64 {
	simplified := make([][][][]float64, len(coordinates))
	for i, multipolygon := range coordinates {
		orbMultipolygon := make(orb.Polygon, len(multipolygon))
		for j, ring := range multipolygon {
			orbRing := make(orb.Ring, len(ring))
			for k, point := range ring {
				orbRing[k] = orb.Point{point[0], point[1]}
			}
			orbMultipolygon[j] = orbRing
		}

		simplifiedMultipolygon := simplify.DouglasPeucker(tolerance).Simplify(orbMultipolygon).(orb.Polygon)
		simplified[i] = make([][][]float64, len(simplifiedMultipolygon))
		for j, ring := range simplifiedMultipolygon {
			simplified[i][j] = make([][]float64, len(ring))
			for k, point := range ring {
				simplified[i][j][k] = []float64{point.X(), point.Y()}
			}
		}
	}
	return simplified
}
