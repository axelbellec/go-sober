// internal/models/drink_option.go
package models

type DrinkOption struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	SizeValue int     `json:"size_value"`
	SizeUnit  string  `json:"size_unit"`
	ABV       float64 `json:"abv"`
}
