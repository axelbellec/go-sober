package user

import (
	"go-sober/internal/dtos"
	"go-sober/internal/models"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) UpdateUserProfile(userID int64, req dtos.UpdateUserProfileRequest) error {
	profile := &models.UserProfile{
		UserID:   userID,
		WeightKg: req.WeightKg,
		Gender:   req.Gender,
	}
	return s.repo.UpsertUserProfile(userID, profile)
}

func (s *Service) GetUserProfile(userID int64) (*models.UserProfile, error) {
	return s.repo.GetUserProfile(userID)
}
