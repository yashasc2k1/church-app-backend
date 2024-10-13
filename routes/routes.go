package routes

import (
	"church-app-backend/controllers"

	"github.com/gofiber/fiber/v2"
)

// SetupUserRoutes sets up routes related to user operations
func SetupUserRoutes(app *fiber.App) {
	// api := app.Group("/api")   // Grouping under /api
	// user := api.Group("/user") // Routes under /api/user

	// POST /api/user/register -> Register a new user
	// user.Post("/register", controllers.RegisterUser)

	// Grouping under /api
	api := app.Group("/api")

	// Routes under /api/user
	user := api.Group("/user")

	// POST /api/user/register -> Register a new user
	user.Post("/register", controllers.RegisterUser)
	user.Post("/register/verify/otp", controllers.VerifyOTP)

	//====================OTP=========================
	user.Post("/otp/generate", controllers.GenerateOTP)

	//===================USER PROFILE=================
	user.Post("/user-profile/", controllers.HandleCreateUserProfile)
	user.Put("/user-profile/", controllers.HandleUpdateUserProfile)

	//===================LOGIN========================
	user.Post("/login", controllers.HandleUserLogin)

}
