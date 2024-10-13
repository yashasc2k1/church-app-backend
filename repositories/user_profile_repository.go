package repositories

import (
	logger "church-app-backend/logger"
	"church-app-backend/models"
	"database/sql"
	"errors"
	"time"
)

func CreateUserProfile(tx *sql.Tx, profile *models.UserProfile) (int64, error) {
	// Initialize the query and arguments slices
	query := "INSERT INTO user_profile (user_id"
	values := "VALUES (? "
	args := []interface{}{profile.UserID}

	// Append fields if they are provided by the user
	if profile.FullName != "" {
		query += ", full_name"
		values += ", ?"
		args = append(args, profile.FullName)
	}
	if !profile.DOB.IsZero() {
		query += ", date_of_birth"
		values += ", ?"
		args = append(args, profile.DOB)
	}
	if profile.MaritalStatus != "" {
		query += ", marital_status"
		values += ", ?"
		args = append(args, profile.MaritalStatus)
	}
	if !profile.WeddingAnniversary.IsZero() {
		query += ", wedding_anniversary"
		values += ", ?"
		args = append(args, profile.WeddingAnniversary)
	}
	if profile.Gender != "" {
		query += ", gender"
		values += ", ?"
		args = append(args, profile.Gender)
	}
	if profile.Profession != "" {
		query += ", profession"
		values += ", ?"
		args = append(args, profile.Profession)
	}

	// Close the query and values
	query += ") " + values + ")"

	// Execute the query
	result, err := tx.Exec(query, args...)
	if err != nil {
		logger.Log.Error("Error Creating User Profile in Database: ", err)
		return -1, err
	}

	// Get tloggerhe last inserted user_id
	profileID, err := result.LastInsertId()
	if err != nil {
		logger.Log.Error("Error Getting Created User-ID: ", err)
		return -1, err
	}
	return profileID, nil
}

func UpdateUserProfile(tx *sql.Tx, profile *models.UserProfile) error {
	// Initialize the query and arguments slices
	query := "UPDATE user_profile SET"
	args := []interface{}{}

	// Append fields if they are provided by the user
	if profile.FullName != "" {
		query += " full_name = ?,"
		args = append(args, profile.FullName)
	}
	if !profile.DOB.IsZero() {
		query += " date_of_birth = ?,"
		args = append(args, profile.DOB)
	}
	if profile.MaritalStatus != "" {
		query += " marital_status = ?,"
		args = append(args, profile.MaritalStatus)
	}
	if !profile.WeddingAnniversary.IsZero() {
		query += " wedding_anniversary = ?,"
		args = append(args, profile.WeddingAnniversary)
	}
	if profile.Gender != "" {
		query += " gender = ?,"
		args = append(args, profile.Gender)
	}
	if profile.Profession != "" {
		query += " profession = ?,"
		args = append(args, profile.Profession)
	}

	// Remove trailing comma
	if len(args) > 0 {
		query = query[:len(query)-1]
	} else {
		// If no fields were provided to update, return an error
		return errors.New("No Fields to update")
	}

	// Add WHERE clause to specify the user
	query += " WHERE user_id = ?"
	args = append(args, profile.UserID)

	// Execute the query
	_, err := tx.Exec(query, args...)
	return err
}

func GetUserProfileByID(tx *sql.Tx, userID int64) (*models.UserProfile, error) {
	var user models.UserProfile

	// Define sql.Null types to handle nullable columns
	var fullName sql.NullString
	var dob sql.NullTime
	var maritalStatus sql.NullString
	var weddingAnniversary sql.NullTime
	var gender sql.NullString
	var profession sql.NullString

	query := `
        SELECT profile_id, user_id, full_name, date_of_birth, marital_status, wedding_anniversary, gender, profession
        FROM user_profile WHERE user_id = ?
    `

	row := tx.QueryRow(query, userID)
	if err := row.Scan(
		&user.ProfileID,
		&user.UserID,
		&fullName,           // Nullable full name
		&dob,                // Nullable date of birth
		&maritalStatus,      // Nullable marital status
		&weddingAnniversary, // Nullable wedding anniversary
		&gender,             // Nullable gender
		&profession,         // Nullable profession
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Log.Error("User Not Found")
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	// Convert sql.Null types to regular types
	user.FullName = fullName.String
	if fullName.Valid {
		user.FullName = fullName.String
	} else {
		user.FullName = "" // Set default or handle as needed
	}

	if dob.Valid {
		user.DOB = dob.Time
	} else {
		user.DOB = time.Time{} // Set default or handle as needed
	}

	if maritalStatus.Valid {
		user.MaritalStatus = maritalStatus.String
	} else {
		user.MaritalStatus = "" // Set default or handle as needed
	}

	if weddingAnniversary.Valid {
		user.WeddingAnniversary = weddingAnniversary.Time
	} else {
		user.WeddingAnniversary = time.Time{} // Set default or handle as needed
	}

	if gender.Valid {
		user.Gender = gender.String
	} else {
		user.Gender = "" // Set default or handle as needed
	}

	if profession.Valid {
		user.Profession = profession.String
	} else {
		user.Profession = "" // Set default or handle as needed
	}

	return &user, nil
}
