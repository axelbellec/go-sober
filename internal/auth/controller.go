package auth

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

// @Summary Sign up a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dtos.UserSignupRequest true "User signup request"
// @Success 201 {object} dtos.UserSignupResponse
// @Failure 400 {object} dtos.ClientError
// @Failure 409 {object} dtos.ClientError
// @Router /auth/signup [post]
func (c *Controller) SignUp(w http.ResponseWriter, r *http.Request) {
	var user dtos.UserSignupRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if user.Email == "" || user.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Check if user already exists
	if _, err := c.service.repo.GetUserByEmail(user.Email); err == nil {
		response := dtos.UserSignupResponse{
			Message: "User already exists",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := c.service.repo.CreateUser(user.Email, user.Password); err != nil {
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	response := dtos.UserSignupResponse{
		Message: "User created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// @Summary Login a user
// @Description Authenticate a user and generate a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dtos.UserLoginRequest true "User login request"
// @Success 200 {object} dtos.UserLoginResponse
// @Failure 400 {object} dtos.ClientError
// @Failure 401 {object} dtos.ClientError
// @Router /auth/login [post]
func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	var credentials dtos.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := c.service.AuthenticateUser(credentials.Email, credentials.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := c.service.GenerateToken(user)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	response := dtos.UserLoginResponse{
		Message: "Login successful",
		Token:   token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// @Summary Get the current user
// @Description Retrieve the current user's information
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} dtos.UserMeResponse
// @Router /auth/me [get]
func (c *Controller) Me(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(constants.UserContextKey).(*models.Claims)

	response := dtos.UserMeResponse{
		UserID: claims.UserID,
		Email:  claims.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
