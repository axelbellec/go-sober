package user

import (
	"encoding/json"
	"net/http"

	"go-sober/internal/constants"
	"go-sober/internal/dtos"
	"go-sober/internal/models"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{service: service}
}

// @Summary Update user profile
// @Description Update the current user's profile information
// @Tags users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param profile body dtos.UpdateUserProfileRequest true "User profile"
// @Success 200 {object} dtos.UserProfileResponse
// @Failure 400 {object} dtos.ClientError
// @Router /users/profile [put]
func (c *Controller) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(constants.UserContextKey).(*models.Claims)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req dtos.UpdateUserProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := c.service.UpdateUserProfile(claims.UserID, req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	profile, err := c.service.GetUserProfile(claims.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if profile == nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	response := dtos.UserProfileResponse{
		WeightKg:  profile.WeightKg,
		Gender:    profile.Gender,
		CreatedAt: profile.CreatedAt,
		UpdatedAt: profile.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// @Summary Get user profile
// @Description Get the current user's profile information
// @Tags users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} dtos.UserProfileResponse
// @Failure 404 {object} dtos.ClientError
// @Router /users/profile [get]
func (c *Controller) GetProfile(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(constants.UserContextKey).(*models.Claims)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	profile, err := c.service.GetUserProfile(claims.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if profile == nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	response := dtos.UserProfileResponse{
		WeightKg:  profile.WeightKg,
		Gender:    profile.Gender,
		UpdatedAt: profile.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
