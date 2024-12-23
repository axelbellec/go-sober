// internal/drinks/service.go
package drinks

import (
	"time"

	"go-sober/internal/models"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetDrinkOptions() ([]models.DrinkOption, error) {
	return s.repo.GetDrinkOptions()
}

func (s *Service) GetDrinkOption(id int) (*models.DrinkOption, error) {
	return s.repo.GetDrinkOption(id)
}

func (s *Service) UpdateDrinkOption(id int, option *models.DrinkOption) error {
	return s.repo.UpdateDrinkOption(id, option)
}

func (s *Service) DeleteDrinkOption(id int) error {
	return s.repo.DeleteDrinkOption(id)
}

func (s *Service) CreateDrinkLog(userID int64, drinkOptionID int64, loggedAt *time.Time) (int64, error) {
	return s.repo.CreateDrinkLog(userID, drinkOptionID, loggedAt)
}

func (s *Service) GetDrinkLogs(userID int64) ([]models.DrinkLog, error) {
	return s.repo.GetDrinkLogs(userID)
}
