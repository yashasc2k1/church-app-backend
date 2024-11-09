package routes

import (
	"church-app-backend/controllers"
	"church-app-backend/middleware"

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

	//===================LOGIN========================
	user.Post("/login", controllers.HandleUserLogin)

	//====================PASSWORD RESET==========================
	user.Post("/forgot-password", controllers.HandleForgotPassword)
	user.Post("/reset-password", controllers.HandleResetPassword)

	//=========================================================================USER MIDDLEWARE INITIALIZED======================================================================
	user.Use(middleware.AuthMiddleware)

	//===================USER PROFILE=================
	user.Post("/user-profile/", controllers.HandleCreateUserProfile)
	user.Put("/user-profile/", controllers.HandleUpdateUserProfile)

	//====================USER RELATED DONATION==========================
	user.Get("/donation-user-list", controllers.HandleGetAllDonationUsers)

	//====================DONATIONS==========================
	user.Post("/donation", controllers.HandleCreateDonation)     // Add a donation
	user.Get("/donation/all", controllers.HandleGetAllDonations) // Get all donations
	user.Get("/donation/total", controllers.HandleTotalDonationCount)
	user.Get("/donation/:userID", controllers.HandleGetDonationByUserID)   // Get donations by user ID
	user.Delete("/donation/:donationID", controllers.HandleDeleteDonation) //Delete donation by donation ID
	user.Put("/donation", controllers.HandleUpdateDonation)                //update donation

}
