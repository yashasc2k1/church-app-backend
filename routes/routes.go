package routes

import (
	"church-app-backend/controllers"

	"github.com/gofiber/fiber/v2"
)

// SetupUserRoutes sets up routes related to user operations
func SetupUserRoutes(app *fiber.App) {
	api := app.Group("/api")   // Grouping under /api
	user := api.Group("/user") // Routes under /api/user

	// POST /api/user/register -> Register a new user
	user.Post("/register", controllers.RegisterUser)
}
