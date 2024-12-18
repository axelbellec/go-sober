package dtos

import "go-sober/internal/models"

type DrinkOptionsResponse struct {
	DrinkOptions []models.DrinkOption `json:"drink_options"`
}

type DrinkOptionResponse struct {
	DrinkOption models.DrinkOption `json:"drink_option"`
}
