package models

import (
	"fmt"
	"time"
)

type DrinkLog struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	DrinkTemplateID int       `json:"drink_template_id"`
	LoggedAt        time.Time `json:"logged_at"`
	DrinkName       string    `json:"drink_name"`
	Type            string    `json:"type"`
	ABV             float64   `json:"abv"`
	SizeValue       int       `json:"size_value"`
	SizeUnit        string    `json:"size_unit"`
}

func (d *DrinkLog) GetVolumeInMl() float64 {
	switch d.SizeUnit {
	case "cl":
		return float64(d.SizeValue) * 10
	case "ml":
		return float64(d.SizeValue)
	default:
		fmt.Printf("Unknown size unit: %s\n", d.SizeUnit)
		return float64(d.SizeValue)
	}
}

func (d *DrinkLog) GetSpecificGravityOfEthanol() float64 {
	// Specific gravity of ethanol in g/ml is 0.789
	return 0.789
}

func (d *DrinkLog) GetAlcoholConsumedInGrams() float64 {
	// Alcohol consumed in grams is equal to :
	// Volume (ml) × ABV × specific gravity of ethanol
	return d.GetVolumeInMl() * d.GetABVInPercent() * d.GetSpecificGravityOfEthanol()
}

func (d *DrinkLog) GetABVInPercent() float64 {
	return d.ABV * 100
}

func (d *DrinkLog) GetStandardDrinks() float64 {
	return d.GetAlcoholConsumedInGrams() / 1000
}
