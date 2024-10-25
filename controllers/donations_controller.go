package controllers

import (
	logger "church-app-backend/logger"
	"church-app-backend/models"
	"church-app-backend/repositories"
	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

/*
curl -X POST http://localhost:6666/api/user/donation \
-H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjk4NzUwMjEsInVzZXJfaWQiOjI0fQ.FC6TyX4X5Uh9zVIN1QJ_0nX0qq7d1b68JPS4B8If1Ag" \
-H "Content-Type: application/json" \
-d '{
"user_id": 16,
"amount": 100,
"purpose": "Charity Donation"
}'
*/
func HandleCreateDonation(c *fiber.Ctx) error {
	// Start a transaction
	tx := c.Locals("tx").(*sql.Tx)

	// Parse Input from request body
	var input models.Donations
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid input")
	}

	insertedID, err := repositories.CreateDonation(tx, &input)
	if err != nil {
		logger.Log.Error("Error Adding new Donation: ", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error generating authentication token")
	}

	logger.Log.Info("Inserted Donation ID: ", insertedID)

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Donation added",
	})
}

/*
curl -X GET http://localhost:6666/api/user/donation/16 -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjk4NzUwMjEsInVzZXJfaWQiOjI0fQ.FC6TyX4X5Uh9zVIN1QJ_0nX0qq7d1b68JPS4B8If1Ag"

*/

func HandleGetDonationByUserID(c *fiber.Ctx) error {
	// Start a transaction
	tx := c.Locals("tx").(*sql.Tx)

	//READ USER ID FROM THE INPUT
	userIDstr := c.Params("userID")

	//CONVERT UserID TO INTEGER
	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		logger.Log.Error("Error Converting User ID into integer")
		return fiber.NewError(fiber.StatusInternalServerError, "Error Converting User ID into integer")
	}

	userDonations, err := repositories.GetDonationsByUserID(tx, int64(userID))
	if err != nil {
		logger.Log.Error("Error Getting Donations for user")
		return fiber.NewError(fiber.StatusInternalServerError, "Error Getting Donations for user")
	}

	// Return the donations in the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"donations": userDonations,
	})
}

/*
curl -X GET http://localhost:6666/api/user/donation/all -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjk4NzUwMjEsInVzZXJfaWQiOjI0fQ.FC6TyX4X5Uh9zVIN1QJ_0nX0qq7d1b68JPS4B8If1Ag"
*/
func HandleGetAllDonations(c *fiber.Ctx) error {
	// Start a transaction
	tx := c.Locals("tx").(*sql.Tx)

	//Get All Donations
	allDonations, err := repositories.GetAllDonations(tx)
	if err != nil {
		logger.Log.Error("Error Getting All Donations")
		return fiber.NewError(fiber.StatusInternalServerError, "Error Getting All Donations")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"donations": allDonations,
	})
}

/*
curl -X GET http://localhost:6666/api/user/donation/total -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjk4NzUwMjEsInVzZXJfaWQiOjI0fQ.FC6TyX4X5Uh9zVIN1QJ_0nX0qq7d1b68JPS4B8If1Ag"
*/
func HandleTotalDonationCount(c *fiber.Ctx) error {
	// Start a transaction
	tx := c.Locals("tx").(*sql.Tx)

	//Get Total Donations
	totalDonations, err := repositories.GetTotalDonationCount(tx)
	if err != nil {
		logger.Log.Error("Error getting total donations")
		return fiber.NewError(fiber.StatusInternalServerError, "Error getting total donations")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"donations": totalDonations,
	})
}
