package repositories

import (
	logger "church-app-backend/logger"
	"church-app-backend/models"
	"database/sql"
	"errors"
	"strings"
)

var UserRepository interface {
	CreateUser(tx *sql.Tx, user *models.User) error
}

// CreateUser inserts a new user into the database.
func CreateUser(tx *sql.Tx, user *models.User) (int64, error) {
	// Start building the query
	query := "INSERT INTO users ("
	values := []interface{}{}
	placeholders := []string{}

	// Check and add phone_number if provided
	if user.PhoneNumber != "" {
		query += "phone_number, "
		values = append(values, user.PhoneNumber)
		placeholders = append(placeholders, "?")
	}

	// Check and add email if provided
	if user.Email != "" {
		query += "email, "
		values = append(values, user.Email)
		placeholders = append(placeholders, "?")
	}

	// Check and add password if provided
	if user.Password != "" {
		query += "password_hash, "
		values = append(values, user.Password)
		placeholders = append(placeholders, "?")
	}

	// Check and add is_verified if provided
	placeholders = append(placeholders, "?")
	query += "is_verified, created_at, updated_at) VALUES ("
	query += strings.Join(placeholders, ", ") + ", NOW(), NOW())"

	// Add the IsVerified value
	values = append(values, user.IsVerified)

	// Execute the query
	result, err := tx.Exec(query, values...)
	if err != nil {
		logger.Log.Error("Error Creating User in Database: ", err)
		return -1, err
	}

	// Get the last inserted user_id
	userID, err := result.LastInsertId()
	if err != nil {
		logger.Log.Error("Error Getting Created User-ID: ", err)
		return -1, err
	}

	user.UserID = userID
	return userID, nil
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
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &user, nil
}

func GetUserByPhoneNumber(tx *sql.Tx, phoneNumber string) (*models.User, error) {
	var user models.User
	query := `
	 	SELECT user_id, phone_number, email, password_hash, is_verified, created_at, updated_at
        FROM users WHERE phone_number = ?
	`

	row := tx.QueryRow(query, phoneNumber)
	if err := row.Scan(&user.UserID, &user.PhoneNumber, &user.Email, &user.Password, &user.IsVerified, &user.CreatedAt, &user.ModifiedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Log.Info("User Not Found")
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &user, nil
}

func GetUserByEmail(tx *sql.Tx, email string) (*models.User, error) {
	var user models.User
	query := `
	 	SELECT user_id, phone_number, email, password_hash, is_verified, created_at, updated_at
        FROM users WHERE email = ?
	`

	row := tx.QueryRow(query, email)
	if err := row.Scan(&user.UserID, &user.PhoneNumber, &user.Email, &user.Password, &user.IsVerified, &user.CreatedAt, &user.ModifiedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Log.Info("User Not Found")
			return nil, sql.ErrNoRows
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
