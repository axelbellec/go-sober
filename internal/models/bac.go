package models

import "time"

// BACStatus represents different levels of blood alcohol content
type BACStatus string

const (
	BACStatusSober       BACStatus = "Sober"
	BACStatusMinimal     BACStatus = "Minimal"
	BACStatusLight       BACStatus = "Light"
	BACStatusMild        BACStatus = "Mild"
	BACStatusSignificant BACStatus = "Significant"
	BACStatusSevere      BACStatus = "Severe"
	BACStatusDangerous   BACStatus = "Dangerous"
)

// GetCategory converts a BACStatus to a simplified BACCategory
func (status BACStatus) GetCategory() BACCategory {
	switch status {
	case BACStatusSober:
		return BACCategorySober
	case BACStatusMinimal, BACStatusLight, BACStatusMild:
		return BACCategoryLight
	default:
		return BACCategoryHeavy
	}
}

// BACCategory represents a simplified BAC status category
type BACCategory string

const (
	BACCategorySober BACCategory = "sober" // BAC == 0
	BACCategoryLight BACCategory = "light" // 0 < BAC < 0.08 (includes minimal, light, mild)
	BACCategoryHeavy BACCategory = "heavy" // BAC >= 0.08 (includes significant, severe, dangerous)
)

type BACCalculationParams struct {
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	WeightKg     float64   `json:"weight_kg"`
	Gender       Gender    `json:"gender" validate:"oneof=male female unknown"`
	TimeStepMins int       `json:"time_step_mins,omitempty"` // Add this field
}

type BACPoint struct {
	Time      time.Time `json:"time"`
	BAC       float64   `json:"bac"`
	Status    BACStatus `json:"status"`
	IsOverBAC bool      `json:"is_over_bac"`
}

// BACSummary provides summary statistics for the BAC calculation
type BACSummary struct {
	MaxBAC             float64   `json:"max_bac"`
	MaxBACTime         time.Time `json:"max_bac_time"`
	SoberSinceTime     time.Time `json:"sober_since_time"`
	TotalDrinks        int       `json:"total_drinks"`
	DrinkingSinceTime  time.Time `json:"drinking_since_time"`
	DurationOverBAC    int       `json:"duration_over_bac"`
	EstimatedSoberTime time.Time `json:"estimated_sober_time"`
}

type BACCalculation struct {
	Timeline []BACPoint `json:"timeline"`
	Summary  BACSummary `json:"summary"`
}
