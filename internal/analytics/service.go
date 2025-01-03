package analytics

import (
	"go-sober/internal/constants"
	"go-sober/internal/dtos"
	"go-sober/internal/models"
	"time"
)

type DrinkStatsRepository interface {
	GetDrinkStats(userID int64, period models.TimePeriod, startDate time.Time, endDate time.Time) ([]models.DrinkStatsPoint, error)
	GetMonthlyBACStats(userID int64, startDate, endDate time.Time) ([]dtos.MonthlyBACStats, error)
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

func (s *Service) GetMonthlyBACStats(userID int64, filters dtos.DrinkStatsFilters) ([]dtos.MonthlyBACStats, error) {
	if filters.StartDate == nil {
		oneYearAgo := time.Now().AddDate(0, -11, 0)
		filters.StartDate = &oneYearAgo
	}
	if filters.EndDate == nil {
		now := time.Now()
		filters.EndDate = &now
	}

	return s.drinkStatsRepo.GetMonthlyBACStats(userID, *filters.StartDate, *filters.EndDate)
}
