package models

import "time"

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // "-" means this won't be included in JSON
	CreatedAt time.Time `json:"created_at"`
}

type Gender string

const (
	Male    Gender = "male"
	Female  Gender = "female"
	Unknown Gender = "unknown"
)
