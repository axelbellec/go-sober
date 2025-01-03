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

// MonthlyBACStats represents BAC category statistics for a specific month
type MonthlyBACStats struct {
	Year   int                        `json:"year"`
	Month  int                        `json:"month"`
	Counts map[models.BACCategory]int `json:"counts"`
	Total  int                        `json:"total"`
}

// MonthlyBACStatsResponse represents the response for monthly BAC statistics
type MonthlyBACStatsResponse struct {
	Stats      []MonthlyBACStats    `json:"stats"`
	Categories []models.BACCategory `json:"categories"`
}
