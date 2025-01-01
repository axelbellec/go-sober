package models

type DrinkTemplate struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	SizeValue int     `json:"size_value"`
	SizeUnit  string  `json:"size_unit"`
	ABV       float64 `json:"abv"`
}

type Quantity struct {
	Value int
	Unit  string
}
