package models

type DrinkParsed struct {
	Name          string  `json:"name"`
	Type          string  `json:"type"`
	SizeValue     float64 `json:"size_value"`
	SizeUnit      string  `json:"size_unit"`
	ABV           float64 `json:"abv"`
	Success       bool    `json:"success"`
	Confidence    float64 `json:"confidence"`
	ErrorMessage  string  `json:"error_message"`
	OriginalInput string  `json:"original_input"`
}
