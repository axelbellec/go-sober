package dtos

import "go-sober/internal/models"

type DrinkOptionsResponse struct {
	DrinkOptions []models.DrinkOption `json:"drink_options"`
}

type DrinkOptionResponse struct {
	DrinkOption models.DrinkOption `json:"drink_option"`
}

// Add after existing DTOs

type UpdateDrinkOptionRequest struct {
	Name      string  `json:"name" validate:"required"`
	Type      string  `json:"type" validate:"required"`
	SizeValue int     `json:"size_value" validate:"required,gt=0"`
	SizeUnit  string  `json:"size_unit" validate:"required"`
	ABV       float64 `json:"abv" validate:"required,gt=0"`
}
