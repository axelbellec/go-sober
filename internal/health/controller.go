package health

import (
	"encoding/json"
	"go-sober/internal/models"
	"net/http"
)

type Controller struct{}

func NewController() *Controller {
	return &Controller{}
}

// @Summary Health check endpoint
// @Description Get API health status
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} models.Health
// @Router /health [get]
func (c *Controller) Health(w http.ResponseWriter, r *http.Request) {
	response := models.Health{
		Status: models.HealthStatusOK,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
