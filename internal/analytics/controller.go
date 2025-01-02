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
	"go-sober/internal/params"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{service: service}
}

// @Summary Get BAC calculation
// @Description Calculate BAC for a user
// @Tags analytics
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param start_time query string true "Start time"
// @Param end_time query string true "End time"
// @Param weight_kg query float64 true "Weight in kg"
// @Param gender query string true "Gender"
// @Param time_step_mins query int true "Time step in minutes"
// @Success 200 {object} dtos.BACCalculationResponse
// @Failure 400 {object} dtos.ClientError
// @Failure 500 {object} dtos.ClientError
// @Router /analytics/timeline/bac [get]
func (c *Controller) GetBAC(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	claims := r.Context().Value(constants.UserContextKey).(*models.Claims)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Parse query parameters
	query := r.URL.Query()

	startTime := params.ParseTimeParam(query.Get("start_time"))
	if startTime == nil {
		http.Error(w, "Invalid start_time parameter", http.StatusBadRequest)
		return
	}

	endTime := params.ParseTimeParam(query.Get("end_time"))
	if endTime == nil {
		http.Error(w, "Invalid end_time parameter", http.StatusBadRequest)
		return
	}

	if endTime.Before(*startTime) {
		http.Error(w, "End time must be after start time", http.StatusBadRequest)
		return
	}

	// Default weight
	weightKg := 70.0
	if weight := query.Get("weight_kg"); weight != "" {
		if _, err := json.Number(weight).Float64(); err != nil {
			http.Error(w, "Invalid weight parameter", http.StatusBadRequest)
			return
		}
	}

	gender := models.ToGender(query.Get("gender"))
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
		StartTime:    *startTime,
		EndTime:      *endTime,
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

// @Summary Get Current BAC
// @Description Get current Blood Alcohol Content for a user
// @Tags analytics
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param weight_kg query float64 true "Weight in kg"
// @Param gender query string true "Gender"
// @Success 200 {object} dtos.CurrentBACResponse
// @Failure 400 {object} dtos.ClientError
// @Failure 500 {object} dtos.ClientError
// @Router /analytics/current/bac [get]
func (c *Controller) GetCurrentBAC(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	claims := r.Context().Value(constants.UserContextKey).(*models.Claims)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse query parameters
	query := r.URL.Query()

	weightKg := 70.0 // Default weight
	if weight := query.Get("weight_kg"); weight != "" {
		if _, err := json.Number(weight).Float64(); err != nil {
			http.Error(w, "Invalid weight parameter", http.StatusBadRequest)
			return
		}
	}

	gender := models.ToGender(query.Get("gender"))
	if gender != models.Male && gender != models.Female && gender != models.Unknown {
		http.Error(w, "Invalid gender parameter", http.StatusBadRequest)
		return
	}

	// Calculate current time range (last 24 hours)
	endTime := time.Now()
	startTime := endTime.Add(-24 * time.Hour)

	req := dtos.BACCalculationRequest{
		StartTime:    startTime,
		EndTime:      endTime,
		WeightKg:     weightKg,
		Gender:       gender,
		TimeStepMins: 1, // Use 1-minute intervals for more precise current BAC
	}

	// Use mapper to convert DTO to model
	calculationParams := mappers.ToBACCalculationParams(req)

	// Calculate BAC points
	bacResults, err := c.service.CalculateBAC(claims.UserID, calculationParams)
	if err != nil {
		http.Error(w, "Error calculating BAC", http.StatusInternalServerError)
		return
	}

	// Get the most recent BAC value
	// by default, the current BAC is zero
	var currentBAC float64 = 0.0
	var currentStatus models.BACStatus
	if len(bacResults.Timeline) > 0 {
		lastPoint := bacResults.Timeline[len(bacResults.Timeline)-1]
		currentBAC = lastPoint.BAC
		currentStatus = lastPoint.Status
	}

	response := dtos.CurrentBACResponse{
		CurrentBAC:         currentBAC,
		BACStatus:          currentStatus,
		LastCalculated:     endTime,
		IsSober:            currentStatus == models.BACStatusSober,
		EstimatedSoberTime: bacResults.Summary.EstimatedSoberTime,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
