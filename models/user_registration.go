package models

type UserRegisterInput struct {
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email,omitempty"` // Optional
	Password    string `json:"password"`
}
