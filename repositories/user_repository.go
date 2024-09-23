package repositories

import (
	logger "church-app-backend/logger"
	"church-app-backend/models"
	"database/sql"
	"errors"
)

var UserRepository interface {
	CreateUser(tx *sql.Tx, user *models.User) error
}

// CreateUser inserts a new user into the database.
func CreateUser(tx *sql.Tx, user *models.User) error {
	query := `
        INSERT INTO users (phone_number, email, password_hash, is_verified, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, NOW(), NOW())
    `

	result, err := tx.Exec(query, user.PhoneNumber, user.Email, user.Password, user.IsVerified)
	if err != nil {
		logger.Log.Error("Error Creating User in Database")
		return err
	}

	// Get the last inserted user_id
	userID, err := result.LastInsertId()
	if err != nil {
		logger.Log.Error("Error Getting Created User-ID")
		return err
	}

	user.UserID = userID
	return nil
}

// GetUserByID fetches a user by their ID.
func GetUserByID(tx *sql.Tx, userID int64) (*models.User, error) {
	var user models.User
	query := `
        SELECT user_id, phone_number, email, password_hash, is_verified, created_at, updated_at
        FROM users WHERE user_id = ?
    `

	row := tx.QueryRow(query, userID)
	if err := row.Scan(&user.UserID, &user.PhoneNumber, &user.Email, &user.Password, &user.IsVerified, &user.CreatedAt, &user.ModifiedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Log.Error("User Not Found")
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// UpdateUser updates the user's details in the database.
func UpdateUser(tx *sql.Tx, user *models.User) error {
	query := `
        UPDATE users SET  phone_number = ?, email = ?, password_hash = ?, is_verified = ?, updated_at = NOW()
        WHERE user_id = ?
    `

	_, err := tx.Exec(query, user.PhoneNumber, user.Email, user.Password, user.IsVerified, user.UserID)
	if err != nil {
		logger.Log.Error("Error updating user")
		return err
	}
	return nil
}

// DeleteUser removes a user by their ID.
func DeleteUser(tx *sql.Tx, userID int64) error {
	query := `DELETE FROM users WHERE user_id = ?`
	_, err := tx.Exec(query, userID)
	if err != nil {
		logger.Log.Error("Error Deleting User")
	}
	return err
}
