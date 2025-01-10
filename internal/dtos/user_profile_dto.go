package dtos

import (
	"go-sober/internal/models"
	"time"
)

type UpdateUserProfileRequest struct {
	WeightKg float64       `json:"weight_kg" validate:"required,gt=0"`
	Gender   models.Gender `json:"gender" validate:"required,oneof=male female unknown"`
}

type UserProfileResponse struct {
	WeightKg  float64       `json:"weight_kg"`
	Gender    models.Gender `json:"gender"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}
