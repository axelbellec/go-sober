package analytics

import (
	"encoding/json"
	"go-sober/internal/constants"
	"go-sober/internal/dtos"
	"go-sober/internal/models"
	"go-sober/internal/params"
	"net/http"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{service: service}
}

// @Summary Get drink statistics
// @Description Get drink statistics for a user
// @Tags analytics
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param period query string true "Time period" Enums(daily, weekly, monthly, yearly)
// @Param start_date query string false "Start date"
// @Param end_date query string false "End date"
// @Success 200 {object} dtos.DrinkStatsResponse
// @Failure 400 {object} dtos.ClientError
// @Failure 500 {object} dtos.ClientError
// @Router /analytics/drink-stats [get]
func (c *Controller) GetDrinkStats(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(constants.UserContextKey).(*models.Claims)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := r.URL.Query()

	period := models.ToTimePeriod(query.Get("period"))
	if period == models.TimePeriodUnknown {
		http.Error(w, "Invalid period parameter", http.StatusBadRequest)
		return
	}

	startDate := params.ParseTimeParam(query.Get("start_date"))
	endDate := params.ParseTimeParam(query.Get("end_date"))

	if startDate != nil && endDate != nil && endDate.Before(*startDate) {
		http.Error(w, "End date must be after start date", http.StatusBadRequest)
		return
	}

	filters := dtos.DrinkStatsFilters{
		Period:    period,
		StartDate: startDate,
		EndDate:   endDate,
	}

	stats, err := c.service.GetDrinkStats(claims.UserID, filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := dtos.DrinkStatsResponse{
		Stats: stats,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// @Summary Get monthly BAC statistics
// @Tags analytics
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param start_date query string false "Start date"
// @Param end_date query string false "End date"
// @Success 200 {object} dtos.MonthlyBACStatsResponse
// @Failure 400 {object} dtos.ClientError
// @Router /analytics/monthly-bac [get]
func (c *Controller) GetMonthlyBACStats(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(constants.UserContextKey).(*models.Claims)
	query := r.URL.Query()

	startDate := params.ParseTimeParam(query.Get("start_date"))
	endDate := params.ParseTimeParam(query.Get("end_date"))

	if startDate != nil && endDate != nil && endDate.Before(*startDate) {
		http.Error(w, "End date must be after start date", http.StatusBadRequest)
		return
	}

	filters := dtos.DrinkStatsFilters{
		StartDate: startDate,
		EndDate:   endDate,
	}

	stats, err := c.service.GetMonthlyBACStats(claims.UserID, filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := dtos.MonthlyBACStatsResponse{
		Stats:      stats,
		Categories: []models.BACCategory{models.BACCategorySober, models.BACCategoryLight, models.BACCategoryHeavy},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
