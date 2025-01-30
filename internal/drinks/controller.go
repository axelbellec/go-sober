// internal/drinks/controller.go
package drinks

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go-sober/internal/constants"
	"go-sober/internal/dtos"
	"go-sober/internal/models"
	"go-sober/internal/params"
	"go-sober/internal/parser"
)

type Controller struct {
	service *Service
	db      *sql.DB
}

func NewController(service *Service, db *sql.DB) *Controller {
	return &Controller{
		service: service,
		db:      db,
	}
}

// @Summary Get all drink templates
// @Description Retrieve all drink templates
// @Tags drinks
// @Accept json
// @Produce json
// @Success 200 {object} dtos.DrinkTemplatesResponse
// @Failure 500 {object} dtos.ClientError
// @Router /drink-templates [get]
func (c *Controller) GetDrinkTemplates(w http.ResponseWriter, r *http.Request) {
	drinkTemplates, err := c.service.GetDrinkTemplates()
	if err != nil {
		http.Error(w, "Could not fetch drink templates", http.StatusInternalServerError)
		return
	}

	response := dtos.DrinkTemplatesResponse{
		DrinkTemplates: drinkTemplates,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// @Summary Get a specific drink template
// @Description Retrieve a specific drink template by ID
// @Tags drinks
// @Accept json
// @Produce json
// @Param id path string true "Drink template ID"
// @Success 200 {object} dtos.DrinkTemplateResponse
// @Failure 404 {object} dtos.ClientError
// @Router /drink-templates/{id} [get]
func (c *Controller) GetDrinkTemplate(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL path
	// The URL pattern "/drink-templates/{id}" needs to be handled with a URL router
	// Since we're using net/http directly, we need to parse the path manually

	// Get the last part of the path
	drinkTemplateID, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	drinkTemplate, err := c.service.GetDrinkTemplate(drinkTemplateID)
	if err != nil {
		http.Error(w, "Drink template not found", http.StatusNotFound)
		return
	}

	response := dtos.DrinkTemplateResponse{
		DrinkTemplate: *drinkTemplate,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// @Summary Create a drink template
// @Description Create a new drink template
// @Tags drinks
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param drinkTemplate body dtos.CreateDrinkTemplateRequest true "New drink template"
// @Success 201 {object} dtos.DrinkTemplateResponse
// @Failure 400 {object} dtos.ClientError
// @Failure 500 {object} dtos.ClientError
// @Router /drink-templates [post]
func (c *Controller) CreateDrinkTemplate(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req dtos.CreateDrinkTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert DTO to model
	drinkTemplate := &models.DrinkTemplate{
		Name:      req.Name,
		Type:      req.Type,
		SizeValue: req.SizeValue,
		SizeUnit:  req.SizeUnit,
		ABV:       req.ABV,
	}

	// Create drink template
	err := c.service.CreateDrinkTemplate(drinkTemplate)
	if err != nil {
		http.Error(w, "Failed to create drink template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// @Summary Update a drink template
// @Description Update a specific drink template by ID
// @Tags drinks
// @Accept json
// @Produce json
// @Param id path string true "Drink template ID"
// @Param drinkTemplate body dtos.UpdateDrinkTemplateRequest true "Updated drink template"
// @Success 204
// @Failure 400 {object} dtos.ClientError
// @Failure 404 {object} dtos.ClientError
// @Failure 500 {object} dtos.ClientError
// @Router /drink-templates/{id} [put]
func (c *Controller) UpdateDrinkTemplate(w http.ResponseWriter, r *http.Request) {

	drinkTemplateID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Parse request body
	var req dtos.UpdateDrinkTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert DTO to model
	drinkTemplate := &models.DrinkTemplate{
		Name:      req.Name,
		Type:      req.Type,
		SizeValue: req.SizeValue,
		SizeUnit:  req.SizeUnit,
		ABV:       req.ABV,
	}

	// Update drink template
	if err := c.service.UpdateDrinkTemplate(drinkTemplateID, drinkTemplate); err != nil {
		if err.Error() == "drink template not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update drink template", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Delete a drink template
// @Description Delete a specific drink template by ID
// @Tags drinks
// @Accept json
// @Produce json
// @Param id path string true "Drink template ID"
// @Success 204
// @Failure 404 {object} dtos.ClientError
// @Failure 500 {object} dtos.ClientError
// @Router /drink-templates/{id} [delete]
func (c *Controller) DeleteDrinkTemplate(w http.ResponseWriter, r *http.Request) {

	drinkTemplateID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Delete drink template
	if err := c.service.DeleteDrinkTemplate(drinkTemplateID); err != nil {
		if err.Error() == "drink template not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete drink template", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Create a drink log
// @Description Create a new drink log for the current user
// @Tags drinks
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param drinkLog body dtos.CreateDrinkLogRequest true "Create drink log request"
// @Success 201 {object} dtos.CreateDrinkLogResponse
// @Failure 400 {object} dtos.ClientError
// @Failure 500 {object} dtos.ClientError
// @Router /drink-logs [post]
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

	// Validate logged_at if provided
	if req.LoggedAt != nil {
		// Ensure logged_at is not in the future
		if req.LoggedAt.After(time.Now()) {
			http.Error(w, "logged_at cannot be in the future", http.StatusBadRequest)
			return
		}
	}

	// Create the drink log and get the ID
	id, err := c.service.CreateDrinkLog(claims.UserID, req)
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

// @Summary Get drink logs for the current user
// @Description Retrieve all drink logs for the current user with optional filters
// @Tags drinks
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param page query int false "Page number (default: 1)" minimum(1)
// @Param page_size query int false "Page size (default: 20, max: 100)" minimum(1) maximum(100)
// @Param start_date query string false "Start date (RFC3339 format)"
// @Param end_date query string false "End date (RFC3339 format)"
// @Param drink_type query string false "Filter by drink type"
// @Param min_abv query number false "Minimum ABV"
// @Param max_abv query number false "Maximum ABV"
// @Param sort_by query string false "Sort by field (logged_at, abv, size_value, name, type)"
// @Param sort_order query string false "Sort order (asc or desc)"
// @Success 200 {object} dtos.GetDrinkLogsResponse
// @Failure 500 {object} dtos.ClientError
// @Router /drink-logs [get]
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

	// Parse pagination parameters
	page, pageSize := params.ParsePaginationParams(r)

	// Parse filter parameters
	filters := dtos.DrinkLogFilters{
		StartDate: params.ParseTimeParam(r.URL.Query().Get("start_date")),
		EndDate:   params.ParseTimeParam(r.URL.Query().Get("end_date")),
		DrinkType: r.URL.Query().Get("drink_type"),
		MinABV:    params.ParseFloatParam(r.URL.Query().Get("min_abv")),
		MaxABV:    params.ParseFloatParam(r.URL.Query().Get("max_abv")),
		SortBy:    r.URL.Query().Get("sort_by"),
		SortOrder: r.URL.Query().Get("sort_order"),
	}

	// Get drink logs from service
	drinkLogs, total, err := c.service.GetDrinkLogs(claims.UserID, page, pageSize, filters)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting drink logs: %v", err), http.StatusInternalServerError)
		return
	}

	response := dtos.GetDrinkLogsResponse{
		DrinkLogs: drinkLogs,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
	}

	// Return drink logs as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// @Summary Parse a drink log
// @Description Parse a drink log and return the drink parsed
// @Tags drinks
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param drinkLog body dtos.ParseDrinkLogRequest true "Parse drink log request"
// @Success 200 {object} dtos.ParseDrinkLogResponse
// @Failure 400 {object} dtos.ClientError
// @Failure 500 {object} dtos.ClientError
// @Router /drink-logs/parse [post]
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

	parser := parser.NewDrinkParser()
	match, err := parser.Parse(req.Text)
	if err != nil {
		http.Error(w, "Could not parse drink description", http.StatusBadRequest)
		return
	}

	response := dtos.ParseDrinkLogResponse{
		DrinkParsed: *match,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// @Summary Update a drink log
// @Description Update a specific drink log for the current user
// @Tags drinks
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Drink log ID"
// @Param drinkLog body dtos.UpdateDrinkLogRequest true "Update drink log request"
// @Success 200 {object} dtos.UpdateDrinkLogResponse
// @Failure 400 {object} dtos.ClientError
// @Failure 401 {object} dtos.ClientError
// @Failure 404 {object} dtos.ClientError
// @Failure 500 {object} dtos.ClientError
// @Router /drink-logs [put]
func (c *Controller) UpdateDrinkLog(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	claims := r.Context().Value(constants.UserContextKey).(*models.Claims)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var req dtos.UpdateDrinkLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate updated_at if provided
	if req.UpdatedAt != nil && req.UpdatedAt.After(time.Now()) {
		http.Error(w, "updated_at cannot be in the future", http.StatusBadRequest)
		return
	}

	// Update the drink log
	err := c.service.UpdateDrinkLog(claims.UserID, req)
	if err != nil {
		http.Error(w, "Failed to update drink log", http.StatusInternalServerError)
		return
	}

	response := dtos.UpdateDrinkLogResponse{
		ID: req.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// @Summary Delete a drink log
// @Description Delete a specific drink log for the current user
// @Tags drinks
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Drink log ID"
// @Success 200 {object} dtos.DeleteDrinkLogResponse
// @Failure 400 {object} dtos.ClientError
// @Failure 401 {object} dtos.ClientError
// @Failure 404 {object} dtos.ClientError
// @Failure 500 {object} dtos.ClientError
// @Router /drink-logs/{id} [delete]
func (c *Controller) DeleteDrinkLog(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	claims := r.Context().Value(constants.UserContextKey).(*models.Claims)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse log ID from path
	logID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid log ID", http.StatusBadRequest)
		return
	}

	// Delete the drink log
	err = c.service.DeleteDrinkLog(claims.UserID, logID)
	if err != nil {
		if err.Error() == "drink log not found or unauthorized" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete drink log: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := dtos.DeleteDrinkLogResponse{
		ID: logID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
