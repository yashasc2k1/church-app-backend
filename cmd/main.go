package main

import (
	"church-app-backend/config"
	logger "church-app-backend/logger"
	"church-app-backend/middleware"
	routes "church-app-backend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v3/log"

	"github.com/joho/godotenv"
)

func main() {
	// Initialize the Fiber app
	app := fiber.New()

	// Initialize the logger
	logger.InitLogger()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Connect to the database
	config.ConnectDatabase()

	// Add the transaction middleware
	app.Use(middleware.TransactionMiddleware(config.DB))

	// Setup routes
	routes.SetupUserRoutes(app)

	// Start the server
	log.Fatal(app.Listen(":6666"))
}
