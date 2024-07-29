package main

import (
	"log"

	"polygon-server/config"
	"polygon-server/database"
	"polygon-server/routes"
	"polygon-server/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load Config
	cfg := config.LoadConfig()

	// Initialize JWT
	utils.InitJWT(cfg.JWTSecret)

	// Connect to MongoDB
	db := database.Connect(cfg)

	// Initialize Validator
	utils.InitValidator()

	// Initialize Fiber
	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())
	app.Use(recover.New())

	// Setup Routes with db dependency
	routes.SetupRoutes(app, db)

	// Start server
	log.Fatal(app.Listen(":" + cfg.Port))
}
