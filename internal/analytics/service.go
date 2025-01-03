package analytics

import (
	"go-sober/internal/constants"
	"go-sober/internal/models"
	"time"
)

type DrinkStatsRepository interface {
	GetDrinkStats(userID int64, period models.TimePeriod, startDate time.Time, endDate time.Time) ([]models.DrinkStats, error)
}

type Service struct {
	drinkStatsRepo DrinkStatsRepository
}

func NewService(drinkStatsRepo DrinkStatsRepository) *Service {
	return &Service{drinkStatsRepo: drinkStatsRepo}
}

func (s *Service) GetDrinkStats(userID int64, period models.TimePeriod, startDate *time.Time, endDate *time.Time) ([]models.DrinkStats, error) {
	if startDate == nil {
		startDate = &constants.DefaultStartDate
	}
	if endDate == nil {
		endDate = &constants.DefaultEndDate
	}

	return s.drinkStatsRepo.GetDrinkStats(userID, period, *startDate, *endDate)
}
