package main

import (
	"church-app-backend/config"
	routes "church-app-backend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v3/log"
)

func main() {
	// Initialize the Fiber app
	app := fiber.New()

	// Connect to the database
	config.ConnectDatabase()

	// Setup routes
	routes.SetupUserRoutes(app)

	// Start the server
	log.Fatal(app.Listen(":6666"))
}
