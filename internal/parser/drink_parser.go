package parser

import (
	"database/sql"
	"fmt"
	"go-sober/internal/embedding"
	"go-sober/internal/models"
	"log"
	"log/slog"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// DrinkMatch represents a matched drink option with its confidence score
type DrinkMatch struct {
	DrinkOption models.DrinkOption
	Confidence  float64
}

// DrinkParser handles the parsing and matching of drink descriptions
type DrinkParser struct {
	drinkOptions     []models.DrinkOption
	drinkEmbeddings  [][]float64
	embeddingService embedding.EmbeddingService
	repository       *embedding.Repository
}

// NewDrinkParser creates a new DrinkParser instance with the given drink options and services
func NewDrinkParser(drinkOptions []models.DrinkOption, embeddingService embedding.EmbeddingService, db *sql.DB) *DrinkParser {
	parser := &DrinkParser{
		drinkOptions:     drinkOptions,
		embeddingService: embeddingService,
		repository:       embedding.NewRepository(db),
	}

	parser.loadOrComputeEmbeddings()
	return parser
}

// loadOrComputeEmbeddings loads existing embeddings from the database or computes new ones
// for each drink option if they don't exist
func (p *DrinkParser) loadOrComputeEmbeddings() {
	p.drinkEmbeddings = make([][]float64, len(p.drinkOptions))
	for i, option := range p.drinkOptions {
		// Try to load from database first
		embedding, err := p.repository.GetEmbedding(option.ID)
		if err == nil {
			p.drinkEmbeddings[i] = embedding
			continue
		}

		// If not found in database, compute and store
		description := p.formatDrinkDescription(option)

		fmt.Println("Getting embedding for", description)
		embedding, err = p.embeddingService.GetEmbedding(description)
		if err != nil {
			log.Printf("Error computing embedding for drink option %d: %v", option.ID, err)
			continue
		}

		// Store in database
		fmt.Println("Storing embedding for", description)
		err = p.repository.StoreEmbedding(option.ID, embedding)
		if err != nil {
			log.Printf("Error storing embedding for drink option %d: %v", option.ID, err)
		}

		p.drinkEmbeddings[i] = embedding
	}
}

// formatDrinkDescription creates a standardized string description of a drink option
func (p *DrinkParser) formatDrinkDescription(option models.DrinkOption) string {
	return fmt.Sprintf("%s %s %d%s", option.Name, option.Type, option.SizeValue, option.SizeUnit)
}

// extractABV attempts to extract an alcohol by volume value from a text string
// Returns the ABV as a decimal (e.g., 0.045 for 4.5%) or an error if no ABV is found
func (p *DrinkParser) extractABV(text string) (float64, error) {

	// Remove spaces and convert to lowercase
	text = strings.ToLower(strings.TrimSpace(text))

	// Match patterns like:
	// - 4.6%, 5.2°, 4.3% ABV, 8%, 0.05
	re := regexp.MustCompile(`(\d+\.?\d*)[%°]|\b0?\.\d+\b`)
	matches := re.FindStringSubmatch(text)

	if len(matches) == 0 {
		return 0, fmt.Errorf("no ABV found in text")
	}

	// Convert the matched value to float
	value, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing ABV value: %w", err)
	}

	// If the value is less than 1, assume it's already in decimal form
	if value < 1 {
		value = value * 100
	}

	// Convert percentage to decimal
	return value / 100, nil

}

// extractQuantity attempts to extract a quantity (value and unit) from a text string
// Supports units: cl, ml, l, pint
// Returns a Quantity struct or an error if no quantity is found
func (p *DrinkParser) extractQuantity(text string) (models.Quantity, error) {
	// 25cl -> 250 ml
	// 1l -> 100 cl
	// 1 pint -> 500 ml

	re := regexp.MustCompile(`(\d+\.?\d*)\s*(cl|ml|l|pint)\b`)
	matches := re.FindStringSubmatch(text)
	if len(matches) == 0 {
		return models.Quantity{}, fmt.Errorf("no quantity found in text")
	}

	value, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return models.Quantity{}, fmt.Errorf("error parsing quantity value: %w", err)
	}

	return models.Quantity{Value: int(value), Unit: matches[2]}, nil
}

// Parse takes a drink description text and returns the best matching drink option
// It also attempts to extract and update the ABV and quantity if present in the text
func (p *DrinkParser) Parse(text string) (*DrinkMatch, error) {
	text = strings.ToLower(strings.TrimSpace(text))

	inputEmbedding, err := p.embeddingService.GetEmbedding(text)
	if err != nil {
		return nil, fmt.Errorf("error computing embedding: %w", err)
	}

	bestMatch := p.findBestMatch(inputEmbedding)
	if bestMatch == nil {
		return nil, fmt.Errorf("could not match drink description")
	}

	// Try to also parse the ABV using regex
	abv, err := p.extractABV(text)
	if err != nil {
		slog.Warn("Could not extract ABV", "text", text)
	} else {
		slog.Debug("Extracted ABV", "text", text, "abvBefore", abv, "abvAfter", bestMatch.DrinkOption.ABV)
		bestMatch.DrinkOption.ABV = abv
	}

	// Try to also parse the size and unit
	quantity, err := p.extractQuantity(text)
	if err != nil {
		fmt.Println(err)
		slog.Warn("Could not extract size and unit", "text", text)
	} else {
		slog.Debug("Extracted size and unit", "text", text, "quantity", quantity.Value, "unit", quantity.Unit)
		bestMatch.DrinkOption.SizeValue = quantity.Value
		bestMatch.DrinkOption.SizeUnit = quantity.Unit
	}

	return bestMatch, nil
}

// findBestMatch finds the drink option with the highest cosine similarity to the input embedding
func (p *DrinkParser) findBestMatch(inputEmbedding []float64) *DrinkMatch {
	if len(p.drinkOptions) == 0 {
		return nil
	}

	// Pre-allocate bestMatch to avoid nil checks
	bestMatch := &DrinkMatch{
		DrinkOption: p.drinkOptions[0],
		Confidence:  cosineSimilarity(inputEmbedding, p.drinkEmbeddings[0]),
	}

	// Start from index 1 since we've already processed index 0
	for i := 1; i < len(p.drinkOptions); i++ {
		similarity := cosineSimilarity(inputEmbedding, p.drinkEmbeddings[i])
		if similarity > bestMatch.Confidence {
			bestMatch.DrinkOption = p.drinkOptions[i]
			bestMatch.Confidence = similarity
		}
	}

	return bestMatch
}

// cosineSimilarity calculates the cosine similarity between two vectors
// Returns a value between -1 and 1, where 1 indicates identical vectors
func cosineSimilarity(a, b []float64) float64 {
	var dotProduct, normA, normB float64

	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}
