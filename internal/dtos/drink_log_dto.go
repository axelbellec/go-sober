package dtos

import (
	"fmt"
	"go-sober/internal/models"
	"time"
)

type CreateDrinkLogRequest struct {
	Name      string     `json:"name" validate:"required"`
	Type      string     `json:"type" validate:"required"`
	SizeValue int        `json:"size_value" validate:"required,gt=0"`
	SizeUnit  string     `json:"size_unit" validate:"required"`
	ABV       float64    `json:"abv" validate:"required,gt=0"`
	LoggedAt  *time.Time `json:"logged_at,omitempty"`
}

type UpdateDrinkLogRequest struct {
	ID        int64      `json:"id" validate:"required"`
	Name      string     `json:"name" validate:"required"`
	Type      string     `json:"type" validate:"required"`
	SizeValue int        `json:"size_value" validate:"required,gt=0"`
	SizeUnit  string     `json:"size_unit" validate:"required"`
	ABV       float64    `json:"abv" validate:"required,gt=0"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type CreateDrinkLogResponse struct {
	ID int64 `json:"id"`
}

type GetDrinkLogsRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type GetDrinkLogsResponse struct {
	DrinkLogs []models.DrinkLog `json:"drink_logs"`
	Total     int               `json:"total"`
	Page      int               `json:"page"`
	PageSize  int               `json:"page_size"`
}

type ParseDrinkLogRequest struct {
	Text string `json:"text"`
}

type ParseDrinkLogResponse struct {
	DrinkParsed models.DrinkParsed `json:"drink_parsed"`
}

type UpdateDrinkLogResponse struct {
	ID int64 `json:"id"`
}

type DeleteDrinkLogResponse struct {
	ID int64 `json:"id"`
}

type DrinkLogFilters struct {
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	DrinkType string     `json:"drink_type"`
	MinABV    *float64   `json:"min_abv"`
	MaxABV    *float64   `json:"max_abv"`
	SortBy    string     `json:"sort_by"`    // e.g., "logged_at", "abv", "size_value"
	SortOrder string     `json:"sort_order"` // "asc" or "desc"
}

func (d *DrinkLogFilters) String() string {
	return fmt.Sprintf("StartDate: %v, EndDate: %v, DrinkType: %v, MinABV: %v, MaxABV: %v, SortBy: %v, SortOrder: %v", d.StartDate, d.EndDate, d.DrinkType, d.MinABV, d.MaxABV, d.SortBy, d.SortOrder)
}
