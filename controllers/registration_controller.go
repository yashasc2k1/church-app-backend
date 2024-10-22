package controllers

import (
	logger "church-app-backend/logger"
	"church-app-backend/models"
	"church-app-backend/repositories"
	"church-app-backend/utils"
	"database/sql"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RegisterUser(c *fiber.Ctx) error {
	tx := c.Locals("tx").(*sql.Tx)
	var input models.UserRegisterInput

	//Check if input is valid
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Input")
	}

	//Check if phone number is not empty
	if len(input.PhoneNumber) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Phone Number is required")
	}

	//Validate if the phone number already exists
	userExists, err := repositories.GetUserByPhoneNumber(tx, input.PhoneNumber)
	if err != sql.ErrNoRows {
		if userExists.IsVerified {
			logger.Log.Error("Error checking user's existence: ", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Error checking user existence")
		}
	}

	userExists, err = repositories.GetUserByEmail(tx, input.Email)
	if err != sql.ErrNoRows {
		if userExists.IsVerified {
			logger.Log.Error("Error checking user's existence: ", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Error checking user existence")
		}
	}
	userID := 0
	if userExists != nil {
		userID = int(userExists.UserID)
	} else {

		// Hash the password
		hashedPassword, err := utils.HashPassword(input.Password)
		if err != nil {
			logger.Log.Error("Error hashing password: ", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Error hashing password")
		}

		//Insert the new user into the database with unverified status
		newUser := models.User{
			PhoneNumber: input.PhoneNumber,
			Email:       input.Email,
			Password:    hashedPassword,
		}
		newUserID, err := repositories.CreateUser(tx, &newUser)
		if err != nil {
			logger.Log.Error("Error creating new User: ", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Error hashing password")
		}
		userID = int(newUserID)
	}

	otpCode := utils.GenerateOTP()

	// Store OTP in the database for both email and phone
	currentTime := time.Now()
	expirationTime := currentTime.Add(5 * time.Minute)

	if err := repositories.CreateOTP(tx, int64(userID), otpCode, expirationTime, false); err != nil {
		logger.Log.Error("Error storing OTP: ", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error storing OTP")
	}

	//Send OTP to Email
	if err := utils.SendEmail(input.Email, "OTP Verification", fmt.Sprintf("Your OTP Verification Code: %s", otpCode)); err != nil {
		logger.Log.Error("Error sending OTP to email: ", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error sending OTP to email")
	}

	// Respond with success message
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User registered successfully. A confirmation email has been sent.",
		"user_id": userID,
	})
}

func VerifyOTP(c *fiber.Ctx) error {
	tx := c.Locals("tx").(*sql.Tx)

	// Check if input is valid
	var input models.VerifyOTP
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Input")
	}

	//Validate OTP
	otp, err := repositories.GetOTP(tx, int64(input.UserID), input.OTPCode)
	if err != nil {
		logger.Log.Error("Error Getting OTP from db: ", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error Getting OTP from Database")
	}

	if otp == nil || otp.IsUsed || otp.ExpiresAt.Before(time.Now()) {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid OTP")
	}

	//Mark OTP as used
	err = repositories.MarkOTPAsUsed(tx, int64(input.UserID))
	if err != nil {
		logger.Log.Error("Error Marking OTP as used: ", err)
		return fiber.NewError(fiber.StatusUnauthorized, "Error Marking OTP as used")
	}

	// Update user verification status in the database

	//Get user
	user, err := repositories.GetUserByID(tx, int64(input.UserID))
	if err != nil {
		logger.Log.Error("Error Retreiving User: ", err)
		return fiber.NewError(fiber.StatusUnauthorized, "Error Retreiving User")
	}

	user.IsVerified = true

	err = repositories.UpdateUser(tx, user)
	if err != nil {
		logger.Log.Error("Error Updating User: ", err)
		return fiber.NewError(fiber.StatusUnauthorized, "Error Updating User")
	}

	// Send a confirmation email after successful verification
	subject := "Welcome to the Church App"
	body := "Dear User, You have successfully registered to our church application"
	err = utils.SendEmail(user.Email, subject, body)
	if err != nil {
		logger.Log.Error("Error sending confirmation email: ", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error sending confirmation email")
	}

	// Respond with success message
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User verified successfully. A confirmation email has been sent.",
	})
}
