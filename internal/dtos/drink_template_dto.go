package dtos

import "go-sober/internal/models"

type DrinkTemplatesResponse struct {
	DrinkTemplates []models.DrinkTemplate `json:"drink_templates"`
}

type DrinkTemplateResponse struct {
	DrinkTemplate models.DrinkTemplate `json:"drink_template"`
}

type UpdateDrinkTemplateRequest struct {
	Name      string  `json:"name" validate:"required"`
	Type      string  `json:"type" validate:"required"`
	SizeValue int     `json:"size_value" validate:"required,gt=0"`
	SizeUnit  string  `json:"size_unit" validate:"required"`
	ABV       float64 `json:"abv" validate:"required,gt=0"`
}

type CreateDrinkTemplateRequest struct {
	Name      string  `json:"name" validate:"required"`
	Type      string  `json:"type" validate:"required"`
	SizeValue int     `json:"size_value" validate:"required,gt=0"`
	SizeUnit  string  `json:"size_unit" validate:"required"`
	ABV       float64 `json:"abv" validate:"required,gt=0"`
}
