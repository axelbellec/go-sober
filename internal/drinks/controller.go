// internal/drinks/controller.go
package drinks

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go-sober/internal/constants"
	"go-sober/internal/dtos"
	"go-sober/internal/embedding"
	"go-sober/internal/models"
	"go-sober/internal/parser"
)

type Controller struct {
	service          *Service
	embeddingService embedding.EmbeddingService
	db               *sql.DB
}

func NewController(service *Service, embeddingService embedding.EmbeddingService, db *sql.DB) *Controller {
	return &Controller{
		service:          service,
		embeddingService: embeddingService,
		db:               db,
	}
}

func (c *Controller) GetDrinkOptions(w http.ResponseWriter, r *http.Request) {
	drinkOptions, err := c.service.GetDrinkOptions()
	if err != nil {
		http.Error(w, "Could not fetch drink options", http.StatusInternalServerError)
		return
	}

	response := dtos.DrinkOptionsResponse{
		DrinkOptions: drinkOptions,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (c *Controller) GetDrinkOption(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL path
	// The URL pattern "/drink-options/{id}" needs to be handled with a URL router
	// Since we're using net/http directly, we need to parse the path manually
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) != 3 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	id := pathParts[2]
	drinkOption, err := c.service.GetDrinkOption(id)
	if err != nil {
		http.Error(w, "Drink option not found", http.StatusNotFound)
		return
	}

	response := dtos.DrinkOptionResponse{
		DrinkOption: *drinkOption,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (c *Controller) CreateDrinkLog(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	claims := r.Context().Value(constants.UserContextKey).(*models.Claims)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var req dtos.CreateDrinkLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate drink_option_id
	if req.DrinkOptionID <= 0 {
		http.Error(w, "Invalid drink_option_id", http.StatusBadRequest)
		return
	}

	// Validate logged_at if provided
	if req.LoggedAt != nil {
		// Ensure logged_at is not in the future
		if req.LoggedAt.After(time.Now()) {
			http.Error(w, "logged_at cannot be in the future", http.StatusBadRequest)
			return
		}
	}

	// Create the drink log and get the ID
	id, err := c.service.CreateDrinkLog(claims.UserID, req.DrinkOptionID, req.LoggedAt)
	if err != nil {
		http.Error(w, "Failed to create drink log: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := dtos.CreateDrinkLogResponse{
		ID: id,
	}

	// Return the created drink log ID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (c *Controller) GetDrinkLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from context (set by auth middleware)
	claims := r.Context().Value(constants.UserContextKey).(*models.Claims)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get drink logs from service
	drinkLogs, err := c.service.GetDrinkLogs(claims.UserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting drink logs: %v", err), http.StatusInternalServerError)
		return
	}

	response := dtos.DrinkLogsResponse{
		DrinkLogs: drinkLogs,
	}

	// Return drink logs as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (c *Controller) ParseDrinkLog(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(constants.UserContextKey).(*models.Claims)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req dtos.ParseDrinkLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	drinkOptions, err := c.service.GetDrinkOptions()
	if err != nil {
		http.Error(w, "Could not fetch drink options", http.StatusInternalServerError)
		return
	}

	parser := parser.NewDrinkParser(drinkOptions, c.embeddingService, c.db)
	match, err := parser.Parse(req.Text)
	if err != nil {
		http.Error(w, "Could not parse drink description", http.StatusBadRequest)
		return
	}

	response := dtos.ParseDrinkLogResponse{
		DrinkOption: match.DrinkOption,
		Confidence:  match.Confidence,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
