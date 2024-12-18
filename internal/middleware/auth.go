package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"go-sober/internal/auth"
	"go-sober/internal/constants"
)

type AuthMiddleware struct {
	service *auth.Service
}

func NewAuthMiddleware(service *auth.Service) *AuthMiddleware {
	return &AuthMiddleware{service: service}
}

func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := extractTokenFromHeader(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, err := m.service.ValidateToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add claims to request context
		ctx := context.WithValue(r.Context(), constants.UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func extractTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header required")
	}

	parts := strings.Fields(authHeader)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("invalid authorization header format")
	}

	return parts[1], nil
}
