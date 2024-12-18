package dtos

import (
	"go-sober/internal/models"
	"time"
)

// BACCalculationRequest represents the input payload for BAC calculation
type BACCalculationRequest struct {
	StartTime    time.Time     `json:"start_time" validate:"required"`
	EndTime      time.Time     `json:"end_time" validate:"required"`
	WeightKg     float64       `json:"weight_kg" validate:"required,gt=0"`
	Gender       models.Gender `json:"gender" validate:"required,oneof=male female unknown"`
	TimeStepMins int           `json:"time_step_mins,omitempty" validate:"omitempty,gt=0"`
}

// BACCalculationResponse represents the output payload for BAC calculation
type BACCalculationResponse struct {
	Timeline []models.BACPoint `json:"timeline"`
	Summary  models.BACSummary `json:"summary,omitempty"`
}
