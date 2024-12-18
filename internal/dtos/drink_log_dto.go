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
