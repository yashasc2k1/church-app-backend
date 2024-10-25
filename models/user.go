package models

import (
	"time"
)

type User struct {
	UserID      int64     `json:"user_id"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
	Password    string    `json:"password_hash"`
	IsVerified  bool      `json:"is_verified"`
	UserType    string    `json:"user_type"`
	CreatedAt   time.Time `json:"created_at"`
	ModifiedAt  time.Time `json:"modified_at"`
}

type UserLoginInput struct {
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type DonationUserList struct {
	UserID      int64  `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}
