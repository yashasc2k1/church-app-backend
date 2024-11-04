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

/*
		curl -X POST http://localhost:6666/api/user/forgot-password \
	     -H "Content-Type: application/json" \
	     -d '{
	           "email": "yashaschandrashek4r@gmail.com"
	         }'
*/
func HandleForgotPassword(c *fiber.Ctx) error {
	tx := c.Locals("tx").(*sql.Tx)
	var input models.ForgotPasswordInput

	// Parse the input
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Input")
	}

	//Get the user by email or phone
	var user *models.User

	if input.Email != "" {
		currUser, err := repositories.GetUserByEmail(tx, input.Email)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Error Retreiving User")
		}

		user = currUser
	} else {
		currUser, err := repositories.GetUserByPhoneNumber(tx, input.PhoneNumber)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Error Retreiving User")
		}

		user = currUser
	}

	//Generate OTP
	otpCode := utils.GenerateOTP()
	expirationTime := time.Now().Add(5 * time.Minute)

	//Store the OTP In the Database
	if err := repositories.CreateOTP(tx, user.UserID, otpCode, expirationTime, false); err != nil {
		logger.Log.Error("Error Inserting OTP into DB")
		return fiber.NewError(fiber.StatusBadRequest, "Error Inserting OTP into DB")
	}

	//Send OTP via Email
	err := utils.SendEmail(user.Email, "Password Reset OTP", fmt.Sprintf("Your OTP code is %s", otpCode))
	if err != nil {
		logger.Log.Error("Error Sending OTP")
		return fiber.NewError(fiber.StatusBadRequest, "Error Sending OTP")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "OTP sent for password reset.",
		"user_id": user.UserID,
	})
}

/*
		curl -X POST http://localhost:6666/api/user/reset-password \
	     -H "Content-Type: application/json" \
	     -d '{
	           "user_id": 3,
	           "otp_code": "482980",
	           "new_password": "helloWorld"
	         }'
*/
func HandleResetPassword(c *fiber.Ctx) error {
	tx := c.Locals("tx").(*sql.Tx)
	var input models.ResetPasswordInput

	// Parse the input
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Input")
	}

	// Validate the OTP
	otp, err := repositories.GetOTP(tx, int64(input.UserID), input.OTPCode)
	if err != nil || otp == nil || otp.IsUsed || otp.ExpiresAt.Before(time.Now()) {
		logger.Log.Error("Invalid or Expired OTP")
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired OTP")
	}

	// Hash the New Password
	hashedPassword, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		logger.Log.Error("Error updating password: ", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error updating password")
	}

	//update the user password
	err = repositories.UpdateUserPassword(tx, hashedPassword, int64(input.UserID))
	if err != nil {
		logger.Log.Error("Error updating new password")
		return fiber.NewError(fiber.StatusUnauthorized, "Error updating new password")

	}

	// Mark OTP as used
	err = repositories.MarkOTPAsUsed(tx, int64(input.UserID))
	if err != nil {
		logger.Log.Error("Error updating OTP: ", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error updating OTP")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password reset successful.",
	})
}
