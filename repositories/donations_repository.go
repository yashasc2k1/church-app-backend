package repositories

import (
	logger "church-app-backend/logger"
	"church-app-backend/models"
	"database/sql"
	"errors"
	"time"
)

var DonationRepository interface {
}

func CreateDonation(tx *sql.Tx, donation *models.Donations) (int64, error) {
	//Check if user exists
	if donation.UserID < 1 {
		return -1, errors.New("User ID not passed")
	}

	_, err := GetUserByID(tx, int64(donation.UserID))
	if err != nil {
		logger.Log.Error("Error with USER ID")
		return -1, err
	}

	query := `INSERT INTO donations (user_id, amount, purpose, donated_at)
	VALUES(?, ?, ?, ?)`

	values := []interface{}{
		donation.UserID, // Foreign key (user_id)
		donation.Amount, // Donation amount
		donation.Purpose,
		time.Now(), // Purpose of the donation
	}

	result, err := tx.Exec(query, values...)
	if err != nil {
		logger.Log.Error("Error Adding new Donation: ", err)
		return -1, err
	}

	donationID, err := result.LastInsertId()
	if err != nil {
		logger.Log.Error("Error Getting Created User-ID: ", err)
		return -1, err
	}

	return donationID, nil
}

func GetAllDonations(tx *sql.Tx) ([]models.Donations, error) {
	// SQL query to select all donations
	query := "SELECT id, user_id, amount, purpose, donated_at FROM donations"

	rows, err := tx.Query(query)
	if err != nil {
		logger.Log.Error("Error Querying Donations: ", err)
		return nil, err
	}
	defer rows.Close()

	// Slice to hold all donations
	var donations []models.Donations

	// Loop through the result set
	for rows.Next() {
		var donation models.Donations
		if err := rows.Scan(&donation.ID, &donation.UserID, &donation.Amount, &donation.Purpose, &donation.DonatedAt); err != nil {
			logger.Log.Error("Error Scanning Donation: ", err)
			return nil, err
		}
		// Append each donation to the slice
		donations = append(donations, donation)
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		logger.Log.Error("Error Iterating Donations: ", err)
		return nil, err
	}

	return donations, nil
}

func GetTotalDonationCount(tx *sql.Tx) (float64, error) {
	// SQL query to sum the total donations
	query := `SELECT SUM(amount) FROM donations`

	var totalDonations float64

	// Execute the query and scan the result into totalDonations
	err := tx.QueryRow(query).Scan(&totalDonations)
	if err != nil {
		return 0, err
	}

	return totalDonations, nil
}

func GetDonationsByUserID(tx *sql.Tx, userID int64) ([]models.Donations, error) {
	// SQL query to select all donations
	query := `SELECT id, user_id, amount, purpose, donated_at 
	FROM donations
	WHERE user_id = ?`

	rows, err := tx.Query(query, userID)
	if err != nil {
		logger.Log.Error("Error Querying Donations for User: ", err)
		return nil, err
	}
	defer rows.Close()

	// Slice to hold all donations
	var donations []models.Donations

	// Loop through the result set
	for rows.Next() {
		var donation models.Donations
		if err := rows.Scan(&donation.ID, &donation.UserID, &donation.Amount, &donation.Purpose, &donation.DonatedAt); err != nil {
			logger.Log.Error("Error Scanning Donation: ", err)
			return nil, err
		}
		// Append each donation to the slice
		donations = append(donations, donation)
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		logger.Log.Error("Error Iterating Donations: ", err)
		return nil, err
	}

	return donations, nil
}
