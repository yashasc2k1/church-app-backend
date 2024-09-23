package models

import (
	"time"

	"gorm.io/gorm"
)

var DB *gorm.DB

type User struct {
	UserID       int64     `json:"user_id"`
	PhoneNumber  string    `json:"phone_number"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	IsVerified   bool      `json:"is_verified"`
	CreatedAt    time.Time `json:"created_at"`
	ModifiedAt   time.Time `json:"modified_at"`
}
