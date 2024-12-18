package auth

import (
	"errors"
	"time"

	"go-sober/internal/models"
	"go-sober/platform"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	repo   *Repository
	config *platform.Config
}

func NewService(repo *Repository, config *platform.Config) *Service {
	return &Service{
		repo:   repo,
		config: config,
	}
}

func (s *Service) GenerateToken(user *models.User) (string, error) {
	claims := &models.Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.Auth.JWT.Secret))
}

func (s *Service) ValidateToken(tokenString string) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.config.Auth.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *Service) AuthenticateUser(email, password string) (*models.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if err := s.repo.ComparePassword(user.Password, password); err != nil {
		return nil, err
	}

	return user, nil
}
