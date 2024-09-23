package models

import "time"

type OTP struct {
	ID        int       `json:"otp_id"`
	UserID    int       `json:"user_id"`
	OTPCode   int       `json:"otp_code"`
	IsUsed    bool      `json:"is_used"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
