package models

import "time"

type UserProfile struct {
	ProfileID          int       `json:"profile_id" bson:"profile_id"`
	UserID             int       `json:"user_id" bson:"user_id"`
	FullName           string    `json:"full_name" bson:"full_name"`
	DOB                time.Time `json:"date_of_birth" bson:"date_of_birth"`
	MaritalStatus      string    `json:"marital_status" bson:"marital_status"`
	WeddingAnniversary time.Time `json:"wedding_anniversary" bson:"wedding_anniversary"`
	Gender             string    `json:"gender" bson:"gender"`
	Profession         string    `json:"profession" bson:"profession"`
}
