package dtos

import (
	"go-sober/internal/models"
	"time"
)

// DrinkStatsResponse represents the response for drink statistics
type DrinkStatsResponse struct {
	Stats []models.DrinkStatsPoint `json:"stats"`
}

// DrinkStatsFilters represents the query parameters for drink statistics
type DrinkStatsFilters struct {
	Period    models.TimePeriod `json:"period" validate:"required"` // daily, weekly, monthly, yearly
	StartDate *time.Time        `json:"start_date,omitempty"`       // Optional start date filter
	EndDate   *time.Time        `json:"end_date,omitempty"`         // Optional end date filter
}
