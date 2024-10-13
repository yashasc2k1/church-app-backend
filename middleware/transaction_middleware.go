package middleware

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
)

// Middleware to start and manage transactions
func TransactionMiddleware(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Start a new transaction
		tx, err := db.Begin()
		if err != nil {
			log.Printf("Error starting transaction: %v", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to start transaction")
		}

		// Store the transaction in the request context
		c.Locals("tx", tx)

		// Call the next handler
		if err := c.Next(); err != nil {
			// Rollback the transaction if an error occurs
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("Error rolling back transaction: %v", rbErr)
			}
			return err
		}

		// Commit the transaction if everything is OK
		if err := tx.Commit(); err != nil {
			log.Printf("Error committing transaction: %v", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to commit transaction")
		}

		return nil
	}
}
