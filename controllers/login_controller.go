package controllers

import (
	logger "church-app-backend/logger"
	models "church-app-backend/models"
	"church-app-backend/repositories"
	"church-app-backend/utils"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func HandleUserLogin(c *fiber.Ctx) error {
	// Start a transaction
	tx := c.Locals("tx").(*sql.Tx)

	// Parse input from request body
	var input models.UserLoginInput
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid input")
	}

	// Ensure at least phone number or email is provided
	if len(input.PhoneNumber) == 0 && len(input.Email) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Either phone number or email is required")
	}

	// Ensure password is provided
	if len(input.Password) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Password is required")
	}

	var user *models.User
	var err error

	// Fetch user details based on phone number or email
	if len(input.PhoneNumber) > 0 {
		user, err = repositories.GetUserByPhoneNumber(tx, input.PhoneNumber)
	} else {
		user, err = repositories.GetUserByEmail(tx, input.Email)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
		}
		logger.Log.Error("Error fetching user: ", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	// Check if the user is verified
	if !user.IsVerified {
		return fiber.NewError(fiber.StatusUnauthorized, "Account not verified. Please complete the registration process")
	}

	// Compare password hash
	if err := utils.ComparePasswords(user.Password, input.Password); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
	}

	// Optionally, generate a token or session (for authorization in future API requests)
	token, err := utils.GenerateJWTToken(uint(user.UserID))
	if err != nil {
		logger.Log.Error("Error generating token: ", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error generating authentication token")
	}

	// Return success response with token
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"user_id": user.UserID,
		"token":   token,
	})
}