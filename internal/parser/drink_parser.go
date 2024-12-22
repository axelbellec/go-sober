package parser

import (
	"database/sql"
	"fmt"
	"go-sober/internal/embedding"
	"go-sober/internal/models"
	"log"
	"math"
	"strings"
)

type DrinkMatch struct {
	DrinkOption models.DrinkOption
	Confidence  float64
}

type DrinkParser struct {
	drinkOptions     []models.DrinkOption
	drinkEmbeddings  [][]float64
	embeddingService embedding.EmbeddingService
	repository       *embedding.Repository
}

func NewDrinkParser(drinkOptions []models.DrinkOption, embeddingService embedding.EmbeddingService, db *sql.DB) *DrinkParser {
	parser := &DrinkParser{
		drinkOptions:     drinkOptions,
		embeddingService: embeddingService,
		repository:       embedding.NewRepository(db),
	}

	parser.loadOrComputeEmbeddings()
	return parser
}

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

func (p *DrinkParser) formatDrinkDescription(option models.DrinkOption) string {
	return fmt.Sprintf("%s %s %d%s", option.Name, option.Type, option.SizeValue, option.SizeUnit)
}

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

	return bestMatch, nil
}

func (p *DrinkParser) findBestMatch(inputEmbedding []float64) *DrinkMatch {
	var bestMatch *DrinkMatch
	highestSimilarity := 0.0

	for i, option := range p.drinkOptions {
		similarity := cosineSimilarity(inputEmbedding, p.drinkEmbeddings[i])
		if similarity > highestSimilarity {
			highestSimilarity = similarity
			bestMatch = &DrinkMatch{
				DrinkOption: option,
				Confidence:  similarity,
			}
		}
	}

	return bestMatch
}

// cosineSimilarity calculates similarity between two vectors
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
