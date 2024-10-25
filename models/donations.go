package models

import "time"

type Donations struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Amount    float64   `json:"amount"`
	Purpose   string    `json:"purpose"`
	DonatedAt time.Time `json:"donated_at"`
}
