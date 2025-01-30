package parser

import (
	"encoding/json"
	"fmt"
	"go-sober/internal/llm"
	"go-sober/internal/models"
	"strconv"
	"strings"
)

// LLMResponse represents the JSON response from the LLM
type LLMResponse struct {
	Beverages []struct {
		Name            string `json:"name"`
		ContainerVolume string `json:"container_volume_value"`
		ContainerUnit   string `json:"container_volume_unit"`
		ContainerType   string `json:"container_type,omitempty"`
		AlcoholContent  string `json:"alcohol_content"`
		Quantity        int    `json:"quantity"`
		Type            string `json:"type"`
	} `json:"beverages"`
}

// DrinkParser handles the parsing of drink descriptions using LLM
type DrinkParser struct{}

// NewDrinkParser creates a new DrinkParser instance
func NewDrinkParser() *DrinkParser {
	return &DrinkParser{}
}

// Parse takes a drink description text and returns the parsed drink details
func (p *DrinkParser) Parse(text string) (*models.DrinkParsed, error) {
	result := &models.DrinkParsed{
		Success:       false,
		OriginalInput: text,
		Confidence:    0,
	}

	// Format the input text as JSON for the LLM
	input := fmt.Sprintf(`{"text": %q}`, text)

	// Call the LLM with the beverage parser prompt
	llmResult, err := llm.Call(BEVERAGE_PARSER_PROMPT, input)
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("error calling LLM: %v", err)
		return result, nil
	}

	// Parse the LLM response
	var response LLMResponse
	if err := json.Unmarshal([]byte(llmResult), &response); err != nil {
		result.ErrorMessage = fmt.Sprintf("error parsing LLM response: %v", err)
		return result, nil
	}

	if len(response.Beverages) == 0 {
		result.ErrorMessage = "no beverages found in the text"
		return result, nil
	}

	// Convert the first beverage to DrinkTemplate
	beverage := response.Beverages[0]

	// Parse alcohol content
	var abv float64 = -1
	if beverage.AlcoholContent != "-1" {
		abvStr := strings.TrimSuffix(beverage.AlcoholContent, "%")
		if parsed, err := strconv.ParseFloat(abvStr, 64); err == nil {
			abv = parsed / 100 // Convert percentage to decimal
		}
	}

	// Convert volume to float64
	var volume float64
	volStr := strings.TrimSpace(beverage.ContainerVolume)
	if parsed, err := strconv.ParseFloat(volStr, 64); err == nil {
		volume = parsed
	}

	// Set the parsed values
	result.Name = beverage.Name
	result.Type = beverage.Type
	result.SizeValue = volume
	result.SizeUnit = beverage.ContainerUnit
	result.ABV = abv
	result.Success = true
	result.Confidence = 1.0

	// Validate the parsed values
	if result.Type == "" && (result.SizeValue <= 0 && result.SizeUnit == "") && result.ABV == -1 {
		result.Success = false
		result.ErrorMessage = "incomplete or invalid drink information"
		result.Confidence = 0.0
	}

	return result, nil
}
