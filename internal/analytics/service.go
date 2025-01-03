package analytics

import (
	"go-sober/internal/constants"
	"go-sober/internal/dtos"
	"go-sober/internal/models"
	"time"
)

type DrinkStatsRepository interface {
	GetDrinkStats(userID int64, period models.TimePeriod, startDate time.Time, endDate time.Time) ([]models.DrinkStatsPoint, error)
}

type Service struct {
	drinkStatsRepo DrinkStatsRepository
}

func NewService(drinkStatsRepo DrinkStatsRepository) *Service {
	return &Service{drinkStatsRepo: drinkStatsRepo}
}

func (s *Service) GetDrinkStats(userID int64, filters dtos.DrinkStatsFilters) ([]models.DrinkStatsPoint, error) {
	if filters.StartDate == nil {
		filters.StartDate = &constants.DefaultStartDate
	}
	if filters.EndDate == nil {
		filters.EndDate = &constants.DefaultEndDate
	}

	return s.drinkStatsRepo.GetDrinkStats(userID, filters.Period, *filters.StartDate, *filters.EndDate)
}
