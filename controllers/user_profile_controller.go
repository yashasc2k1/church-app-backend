package controllers

import (
	logger "church-app-backend/logger"
	"church-app-backend/models"
	"church-app-backend/repositories"
	"database/sql"
	"strconv"

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

/*
	curl -X GET "http://127.0.0.1:6666/api/user/user-profile/2" \
	     -H "Content-Type: application/json" \
	     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzEyNzg3MzEsInVzZXJfaWQiOjV9.f4e7-hzdi5gkJQwnSynxVCJRanftAPwmhrOBaIKSzKc"
*/
func HandleGetUserProfile(c *fiber.Ctx) error {
	tx := c.Locals("tx").(*sql.Tx)
	id := c.Params("userID")

	//CONVERT UserID TO INTEGER
	userID, err := strconv.Atoi(id)
	if err != nil {
		logger.Log.Error("Error converting User ID string to integer")
		return fiber.NewError(fiber.StatusInternalServerError, "Error converting User ID string to integer")
	}

	userProfile, err := repositories.GetUserProfileByID(tx, int64(userID))
	if err != nil {
		logger.Log.Error("Error Getting User Profile")
		return fiber.NewError(fiber.StatusInternalServerError, "Error Getting User Profile")

	}
	// Return the donations in the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user_profile": userProfile,
	})
}
