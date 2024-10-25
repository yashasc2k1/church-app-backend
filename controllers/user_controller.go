package controllers

import (
	logger "church-app-backend/logger"
	"church-app-backend/repositories"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func HandleGetAllDonationUsers(c *fiber.Ctx) error {
	tx := c.Locals("tx").(*sql.Tx)

	allUsers, err := repositories.GetAllDonationUsers(tx)
	if err != nil {
		logger.Log.Error("Error getting all donation user list: ", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error getting all donation user list")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": " Donation Users retreived successfully",
		"users":   allUsers,
	})
}
