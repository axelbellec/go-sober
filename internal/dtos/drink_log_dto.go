package dtos

import (
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
	LoggedAt  *time.Time `json:"logged_at,omitempty"`
}

type CreateDrinkLogResponse struct {
	ID int64 `json:"id"`
}

type DrinkLogsResponse struct {
	DrinkLogs []models.DrinkLog `json:"drink_logs"`
}

type ParseDrinkLogRequest struct {
	Text string `json:"text"`
}

type ParseDrinkLogResponse struct {
	DrinkTemplate models.DrinkTemplate `json:"drink_template"`
	Confidence    float64              `json:"confidence"`
}
