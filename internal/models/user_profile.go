package models

import "time"

type UserProfile struct {
	UserID    int64     `json:"user_id"`
	WeightKg  float64   `json:"weight_kg"`
	Gender    Gender    `json:"gender"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
