package controllers

import (
	logger "church-app-backend/logger"
	"church-app-backend/models"
	"church-app-backend/repositories"
	"church-app-backend/utils"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

/*
curl -X POST http://localhost:6666/api/user/donation \
-H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzA4OTk3NTIsInVzZXJfaWQiOjF9.63a7r7TA1LgsNf0kIcWlolA3xgb0HWAcAa-WO9YtSTE" \
-H "Content-Type: application/json" \
-d '{
"user_id": 3,
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

	input.DonatedAt = time.Now()
	insertedID, err := repositories.CreateDonation(tx, &input)
	if err != nil {
		logger.Log.Error("Error Adding new Donation: ", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error generating authentication token")
	}

	logger.Log.Info("Inserted Donation ID: ", insertedID)

	//===================SEND MAIL TO USER REGARDING DONATION=============================

	//get user info using user-id
	user, err := repositories.GetUserByID(tx, int64(input.UserID))
	if err != nil {
		logger.Log.Error("Error Retreiving User Info")
		return fiber.NewError(fiber.StatusInternalServerError, "Error Retreiving User Info")
	}

	//get user profile
	userProfile, err := repositories.GetUserProfileByID(tx, user.UserID)
	if err != nil {
		logger.Log.Error("Error Retreiving User Profile Info")
		return fiber.NewError(fiber.StatusInternalServerError, "Error Retreiving User Profile Info")
	}

	logger.Log.Info(fmt.Sprintf("User EMAIL: %s", user.Email))

	subject := "Thank You for Your Donation"
	body := fmt.Sprintf(`Dear %s,

We are incredibly grateful for your recent donation of ₹%.2f towards %s. Your generosity and support make a meaningful difference in our mission and help us achieve our goals.
	
Donation Details:
- Amount: ₹%.2f
- Purpose: %s
- Date: %s
	
Thank you for being a valued part of our community and for helping us continue our work. Your contribution has a profound impact, and we deeply appreciate your support.
	
If you have any questions about this donation or our ongoing projects, please feel free to contact us.
	
Warm regards,
Church
	`,
		userProfile.FullName,
		input.Amount,
		input.Purpose,
		input.Amount,
		input.Purpose,
		input.DonatedAt.Format("January 2, 2006"),
	)

	err = utils.SendEmail(user.Email, subject, body)

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
curl -X GET http://localhost:6666/api/user/donation/all -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzEyNzg3MzEsInVzZXJfaWQiOjV9.f4e7-hzdi5gkJQwnSynxVCJRanftAPwmhrOBaIKSzKc"
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

/*
	curl -X PUT http://localhost:6666/api/user/donation \
	  -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzEyMzQ0ODksInVzZXJfaWQiOjF9.yu78LKmsWCQqoXjM3PlPvB0Hc9GoMCBodLRv1dOcwzg" \
	  -H "Content-Type: application/json" \
	  -d '{
		"id": 13,
	    "user_id": 2,
	    "amount": 6000,
	    "purpose": "Education"
	  }'
*/
func HandleUpdateDonation(c *fiber.Ctx) error {
	tx := c.Locals("tx").(*sql.Tx)

	// Parse input from request body
	var input models.Donations
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid input")
	}

	// Retrieve existing donation from the database
	existingDonation, err := repositories.GetDonationByID(tx, int64(input.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Log.Error("Donation Does not exist")
		} else {
			logger.Log.Error("Error getting donation by donation id")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Error retrieving donation")
	}

	//update existing fields based on the input
	if input.UserID != 0 {
		//check if the user exists
		_, err := repositories.GetUserByID(tx, int64(input.UserID))
		if err != nil {
			if err == sql.ErrNoRows {
				logger.Log.Error("User does not exist")
				return fiber.NewError(fiber.StatusInternalServerError, "User does not exist")
			}
			logger.Log.Error("Error getting user from user-id")
			return fiber.NewError(fiber.StatusInternalServerError, "Error getting user from user-id")

		}

		existingDonation.UserID = input.UserID
	}

	if input.Amount != 0 {
		existingDonation.Amount = input.Amount
	}

	if input.Purpose != "" {
		existingDonation.Purpose = input.Purpose
	}

	existingDonation.UpdatedAt = time.Now()

	err = repositories.UpdateDonationByDonationID(tx, *existingDonation)
	if err != nil {
		logger.Log.Error("Error updating Donation")
		return fiber.NewError(fiber.StatusInternalServerError, "Error updating Donation")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Donation updated successfully",
	})
}

/*
		curl -X DELETE http://localhost:6666/api/user/donation/14 \
	  	-H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzEyMzQ0ODksInVzZXJfaWQiOjF9.yu78LKmsWCQqoXjM3PlPvB0Hc9GoMCBodLRv1dOcwzg"
*/
func HandleDeleteDonation(c *fiber.Ctx) error {
	tx := c.Locals("tx").(*sql.Tx)

	// Get donation ID from URL params and convert to integer
	id := c.Params("donationID")
	donationID, err := strconv.Atoi(id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid donation ID")
	}

	// Call repository to delete the donation by ID
	err = repositories.DeleteDonationByDonationID(tx, int64(donationID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error deleting donation")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Donation deleted successfully",
	})
}
