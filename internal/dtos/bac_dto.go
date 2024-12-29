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
	TimeStepMins int           `json:"time_step_mins" validate:"required,gt=0"`
}

// BACCalculationResponse represents the output payload for BAC calculation
type BACCalculationResponse struct {
	Timeline []models.BACPoint `json:"timeline"`
	Summary  models.BACSummary `json:"summary,omitempty"`
}

type CurrentBACResponse struct {
	CurrentBAC         float64          `json:"current_bac"`
	BACStatus          models.BACStatus `json:"bac_status"`
	LastCalculated     time.Time        `json:"last_calculated"`
	IsSober            bool             `json:"is_sober"`
	EstimatedSoberTime time.Time        `json:"estimated_sober_time"`
}
