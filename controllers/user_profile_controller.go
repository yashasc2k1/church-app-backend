package controllers

import (
	logger "church-app-backend/logger"
	"church-app-backend/models"
	"church-app-backend/repositories"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func HandleCreateUserProfile(c *fiber.Ctx) error {
	tx := c.Locals("tx").(*sql.Tx)
	var input models.UserProfile

	// Parse the request body into the UserProfile struct
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid input")
	}

	//Check if User Exists
	user, err := repositories.GetUserByID(tx, int64(input.UserID))
	if err != nil && err != sql.ErrNoRows {
		logger.Log.Error("Error finding User: ", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error finding User")
	}
	if user == nil {
		logger.Log.Error("User does not exist: ", err)
		return fiber.NewError(fiber.StatusInternalServerError, "User does not exist")
	}

	// Insert the new profile into the database
	profileID, err := repositories.CreateUserProfile(tx, &input)
	if err != nil {
		logger.Log.Error("Error creating new User Profile: ", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error creating User Profile")
	}

	// Respond with success
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "User Profile added successfully",
		"profile_id": profileID,
	})
}

func HandleUpdateUserProfile(c *fiber.Ctx) error {
	tx := c.Locals("tx").(*sql.Tx)
	var input models.UserProfile

	// Parse the request body into the UserProfile struct
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid input")
	}

	// Check if profile exists before updating
	_, err := repositories.GetUserProfileByID(tx, int64(input.UserID))
	if err != nil {
		logger.Log.Error("Error retrieving User Profile: ", err)
		return fiber.NewError(fiber.StatusNotFound, "Profile not found")
	}

	// Update the profile in the database
	err = repositories.UpdateUserProfile(tx, &input)
	if err != nil {
		logger.Log.Error("Error updating User Profile: ", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error updating User Profile")
	}

	// Respond with success
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User Profile updated successfully",
	})
}
