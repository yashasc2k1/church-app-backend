package repositories

import (
	logger "church-app-backend/logger"
	"church-app-backend/models"
	"database/sql"
	"time"
)

// Create new OTP
func CreateOTP(tx *sql.Tx, userID int64, otpCode string, expiresAt time.Time, isUsed bool) error {
	query := "INSERT INTO otp_verification (user_id, otp_code,created_at, expires_at, is_used) VALUES (?, ?, ?, ?, ?)"
	_, err := tx.Exec(query, userID, otpCode, time.Now(), expiresAt, isUsed)
	if err != nil {
		logger.Log.Error("Error Inserting OTP into Database")
	}
	return err
}

// Get OTP for user
func GetOTP(tx *sql.Tx, userID int64, otpCode string) (*models.OTP, error) {
	var otp models.OTP
	query := "SELECT otp_id, user_id, otp_code, created_at, expires_at, is_used FROM otp_verification WHERE user_id = ? AND otp_code = ? AND is_used = 0"

	err := tx.QueryRow(query, userID, otpCode).Scan(&otp.ID, &otp.UserID, &otp.OTPCode, &otp.CreatedAt, &otp.ExpiresAt, &otp.IsUsed)
	if err != nil {
		return nil, err
	}

	return &otp, nil
}

// Mark OTP as used
func MarkOTPAsUsed(tx *sql.Tx, userID int64) error {
	query := "UPDATE otp_verification SET is_used = 1 WHERE user_id = ?"
	_, err := tx.Exec(query, userID)
	if err != nil {
		logger.Log.Error("Error Updating OTP is_used status")
	}
	return err
}

// Delete Expired OTP
func DeleteExpiredOTPs(tx *sql.Tx) error {
	query := "DELETE FROM otp_verification WHERE expires_at < NOW()"
	_, err := tx.Exec(query)
	if err != nil {
		logger.Log.Error("Error Deleting Expired OTP")
	}
	return err
}
