package repositories

import (
	logger "church-app-backend/logger"
	"database/sql"
)

func UpdateUserPassword(tx *sql.Tx, passwordHash string, userID int64) error {
	// SQL query to select relevant user data from both tables
	query := `UPDATE users
	SET password_hash = ?
	WHERE user_id = ?`

	_, err := tx.Exec(query, passwordHash, userID)
	if err != nil {
		logger.Log.Error("Error Updating Password in DB")
		return err
	}

	return nil
}
