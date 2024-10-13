package controllers

import (
	logger "church-app-backend/logger"
	"church-app-backend/repositories"
	"church-app-backend/utils"
	"database/sql"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GenerateOTP(c *fiber.Ctx) error {
	tx := c.Locals("tx").(*sql.Tx)
	var input struct {
		UserID int `json:"user_id"`
	}

	// Check if input is valid
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Input")
	}

	// Check if user exists
	user, err := repositories.GetUserByID(tx, int64(input.UserID))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	// Generate new OTP
	otpCode := utils.GenerateOTP() // Implement this function to generate a random OTP

	// Store the new OTP in the database
	currentTime := time.Now()
	expirationTime := currentTime.Add(5 * time.Minute)

	if err := repositories.CreateOTP(tx, user.UserID, otpCode, expirationTime, false); err != nil {
		logger.Log.Error("Error storing OTP: ", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error storing OTP")
	}

	// Send OTP to both email and phone
	if err := utils.SendEmail(user.Email, "OTP Verfication", fmt.Sprintf("Your OTP Verification Code: %s", otpCode)); err != nil {
		logger.Log.Error("Error sending OTP to email: ", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error sending OTP to email")
	}

	// Respond with success message
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "A new OTP has been sent to your phone and email.",
	})
}
