// internal/drinks/service.go
package drinks

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

func (s *Service) GetDrinkTemplates() ([]models.DrinkTemplate, error) {
	return s.repo.GetDrinkTemplates()
}

func (s *Service) GetDrinkTemplate(id int) (*models.DrinkTemplate, error) {
	return s.repo.GetDrinkTemplate(id)
}

func (s *Service) CreateDrinkTemplate(template *models.DrinkTemplate) error {
	return s.repo.CreateDrinkTemplate(template)
}

func (s *Service) UpdateDrinkTemplate(id int, template *models.DrinkTemplate) error {
	return s.repo.UpdateDrinkTemplate(id, template)
}

func (s *Service) DeleteDrinkTemplate(id int) error {
	return s.repo.DeleteDrinkTemplate(id)
}

func (s *Service) CreateDrinkLog(userID int64, createDrinkLogRequest dtos.CreateDrinkLogRequest) (int64, error) {
	return s.repo.CreateDrinkLog(userID, createDrinkLogRequest)
}

func (s *Service) UpdateDrinkLog(userID int64, updateDrinkLogRequest dtos.UpdateDrinkLogRequest) error {
	return s.repo.UpdateDrinkLog(userID, updateDrinkLogRequest)
}

func (s *Service) DeleteDrinkLog(userID int64, logID int64) error {
	return s.repo.DeleteDrinkLog(logID, userID)
}

func (s *Service) GetDrinkLogs(userID int64, page, pageSize int, filters dtos.DrinkLogFilters) ([]models.DrinkLog, int, error) {
	return s.repo.GetDrinkLogs(userID, page, pageSize, filters)
}
