package analytics

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"go-sober/internal/constants"
	"go-sober/internal/dtos"
	"go-sober/internal/mappers"
	"go-sober/internal/models"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) GetBAC(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	claims := r.Context().Value(constants.UserContextKey).(*models.Claims)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Parse query parameters
	query := r.URL.Query()

	startTime, err := time.Parse(time.RFC3339, query.Get("start_time"))
	if err != nil {
		http.Error(w, "Invalid start_time parameter", http.StatusBadRequest)
		return
	}

	endTime, err := time.Parse(time.RFC3339, query.Get("end_time"))
	if err != nil {
		http.Error(w, "Invalid end_time parameter", http.StatusBadRequest)
		return
	}

	weightKg := 70.0 // Default weight
	if weight := query.Get("weight_kg"); weight != "" {
		if _, err := json.Number(weight).Float64(); err != nil {
			http.Error(w, "Invalid weight parameter", http.StatusBadRequest)
			return
		}
	}

	gender := models.Gender(query.Get("gender"))
	if gender != models.Male && gender != models.Female && gender != models.Unknown {
		http.Error(w, "Invalid gender parameter", http.StatusBadRequest)
		return
	}

	timeStepMins := 15 // Default value
	if step := query.Get("time_step_mins"); step != "" {
		if stepInt, err := strconv.Atoi(step); err == nil {
			if stepInt > 0 {
				timeStepMins = stepInt
			}
		}
	}

	req := dtos.BACCalculationRequest{
		StartTime:    startTime,
		EndTime:      endTime,
		WeightKg:     weightKg,
		Gender:       gender,
		TimeStepMins: timeStepMins,
	}

	// Use mapper to convert DTO to model
	calculationParams := mappers.ToBACCalculationParams(req)

	// Calculate BAC points
	bacResults, err := c.service.CalculateBAC(claims.UserID, calculationParams)
	if err != nil {
		http.Error(w, "Error calculating BAC", http.StatusInternalServerError)
		return
	}

	// Use mapper to convert model to response DTO
	response := mappers.ToBACCalculationResponse(bacResults)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
