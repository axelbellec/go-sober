package dtos

import (
	"go-sober/internal/models"
	"time"
)

type CreateDrinkLogRequest struct {
	DrinkOptionID int64      `json:"drink_option_id"`
	LoggedAt      *time.Time `json:"logged_at,omitempty"`
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
	DrinkOption models.DrinkOption `json:"drink_option"`
	Confidence  float64            `json:"confidence"`
}
