package routes

import (
	"polygon-server/handlers"
	"polygon-server/middleware"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(app *fiber.App, db *mongo.Database) {
	userHandler := &handlers.UserHandler{DB: db}

	api := app.Group("/api")
	api.Post("/login", userHandler.Login)

	// Protect routes with JWT middleware
	protected := api.Group("/v1", middleware.JWTMiddleware())
	protected.Get("/users", userHandler.GetUsers)
	protected.Post("/users", userHandler.CreateUser)
}
