package analytics

import (
	"encoding/json"
	"go-sober/internal/constants"
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
// @Success 200 {array} models.DrinkStats
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

	stats, err := c.service.GetDrinkStats(claims.UserID, period, startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}
