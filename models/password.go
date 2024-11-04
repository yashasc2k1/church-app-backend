package models

type ForgotPasswordInput struct {
	Email       string `json:"email" bson:"email"`
	PhoneNumber string `json:"phone_number" bson:"phone_number"`
}

type ResetPasswordInput struct {
	OTPCode     string `json:"otp_code" bson:"otp_code"`
	UserID      int    `json:"user_id" bson:"user_id"`
	NewPassword string `json:"new_password" bson:"new_password"`
}
