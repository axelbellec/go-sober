package bac

import (
	"fmt"
	"math"
	"sort"
	"time"

	"go-sober/internal/constants"
	"go-sober/internal/dtos"
	"go-sober/internal/models"
)

type DrinkLogRepository interface {
	GetDrinkLogs(userID int64, page, pageSize int, filters dtos.DrinkLogFilters) ([]models.DrinkLog, int, error)
}

type Service struct {
	drinkLogRepo DrinkLogRepository
}

// Constants for BAC calculation
const (
	maleWidmarkFactor   = 0.68
	femaleWidmarkFactor = 0.55

	metabolismRatePerHour = 0.015
	metabolismRatePerMin  = metabolismRatePerHour / 60
	absorptionTimeMin     = 60 // 1 hour which is the maximum absorption time
	maxPhysiologicalBAC   = 0.55

	// // Beta distribution parameters
	// alpha = 2.0
	// beta  = 3.0
)

type BACTimeline struct {
	Timeline []models.BACPoint `json:"timeline"`
}

func NewService(drinkLogRepo DrinkLogRepository) *Service {
	return &Service{
		drinkLogRepo: drinkLogRepo,
	}
}

// Update the CalculateBAC method signature and implementation
func (s *Service) CalculateBAC(userID int64, params models.BACCalculationParams) (models.BACCalculation, error) {

	filters := dtos.DrinkLogFilters{
		StartDate: &params.StartTime,
		EndDate:   &params.EndTime,
	}
	drinks, _, err := s.drinkLogRepo.GetDrinkLogs(userID, 1, constants.MaxPageSize, filters)
	if err != nil {
		return models.BACCalculation{}, err
	}

	// Calculate BAC
	timeline, err := s.calculateBAC(drinks, params)
	if err != nil {
		return models.BACCalculation{}, err
	}

	totalDrinksConsumed := len(drinks)

	response := models.BACCalculation{
		Timeline: make([]models.BACPoint, len(timeline.Timeline)),
		Summary:  s.calculateBACSummary(timeline.Timeline, totalDrinksConsumed, params.TimeStepMins),
	}

	// Map timeline points
	for i, point := range timeline.Timeline {
		response.Timeline[i] = models.BACPoint{
			Time:      point.Time,
			BAC:       point.BAC,
			Status:    s.getBACStatus(point.BAC),
			IsOverBAC: point.BAC > 0.08,
		}
	}

	return response, nil
}

func (s *Service) calculateBACPoints(drinks []models.DrinkLog, startTime, endTime time.Time,
	bodyWeightGrams float64, widmarkFactor float64, timeStepMin int) []models.BACPoint {

	var bacPoints []models.BACPoint
	currentTime := startTime
	timeStep := time.Duration(timeStepMin) * time.Minute

	for !currentTime.After(endTime) {
		bac := s.calculateBACAtTime(drinks, currentTime, bodyWeightGrams, widmarkFactor)

		if bac > maxPhysiologicalBAC {
			fmt.Printf("BAC %f is above max physiological BAC of %f at %s\n", bac, maxPhysiologicalBAC, currentTime)
			bac = maxPhysiologicalBAC
		}

		bacPoints = append(bacPoints, models.BACPoint{
			Time: currentTime,
			BAC:  bac,
		})

		currentTime = currentTime.Add(timeStep)
	}

	return bacPoints
}

func (s *Service) calculateBACAtTime(drinks []models.DrinkLog, currentTime time.Time,
	bodyWeightGrams, widmarkFactor float64) float64 {

	var totalBAC float64

	for _, drink := range drinks {
		timeElapsed := currentTime.Sub(drink.LoggedAt).Minutes()
		if timeElapsed < 0 {
			continue
		}

		drinkBAC := s.calculateSingleDrinkBAC(drink, timeElapsed, bodyWeightGrams, widmarkFactor)
		totalBAC += drinkBAC
	}

	return totalBAC
}

func (s *Service) calculateSingleDrinkBAC(drink models.DrinkLog, timeElapsed, bodyWeightGrams,
	widmarkFactor float64) float64 {

	initialBAC := drink.GetAlcoholConsumedInGrams() / (bodyWeightGrams * widmarkFactor)
	absorptionFactor := s.calculateAbsorptionFactor(timeElapsed)

	metabolized := metabolismRatePerMin * timeElapsed
	return math.Max(0, (initialBAC*absorptionFactor)-metabolized)
}

func (s *Service) calculateAbsorptionFactor(timeElapsed float64) float64 {
	if timeElapsed > absorptionTimeMin {
		return 1.0
	}

	// Linear increase from 0 to 1
	normalizedTime := timeElapsed / absorptionTimeMin

	// // Beta approximation
	// // The linear model is less accurate because alcohol
	// // absorption isn't constant - it varies based on many
	// // factors and follows more of an S-curve pattern.
	// // That's why the beta distribution (which creates an S-curve)
	// // is preferred for more accurate BAC calculations.
	// betaAbsorption := s.betaCurve(normalizedTime)
	// return betaAbsorption

	return normalizedTime
}

// func (s *Service) betaCurve(t float64) float64 {
// 	// Normalize time to [0,1] interval
// 	x := math.Max(0, math.Min(1, t))

// 	// Beta distribution formula
// 	numerator := math.Pow(x, alpha-1) * math.Pow(1-x, beta-1)
// 	denominator := 1.0 // B(alpha,beta) normalization constant

// 	// Scale to ensure maximum value is 1.0
// 	maxValue := math.Pow((alpha-1)/(alpha+beta-2), alpha-1) *
// 		math.Pow((beta-1)/(alpha+beta-2), beta-1)

// 	return (numerator / denominator) / maxValue
// }

func (s *Service) getWidmarkFactor(gender models.Gender) float64 {
	switch gender {
	case models.Female:
		return femaleWidmarkFactor
	default:
		return maleWidmarkFactor
	}
}

// calculateBAC handles the core BAC calculation logic
func (s *Service) calculateBAC(drinks []models.DrinkLog, params models.BACCalculationParams) (BACTimeline, error) {

	// Sort drinks chronologically
	sort.Slice(drinks, func(i, j int) bool {
		return drinks[i].LoggedAt.Before(drinks[j].LoggedAt)
	})

	widmarkFactor := s.getWidmarkFactor(params.Gender)
	bodyWeightGrams := params.WeightKg * 1000

	points := s.calculateBACPoints(drinks, params.StartTime, params.EndTime, bodyWeightGrams, widmarkFactor, params.TimeStepMins)
	return BACTimeline{Timeline: points}, nil
}

// calculateBACSummary generates a summary of the BAC timeline
func (s *Service) calculateBACSummary(timeline []models.BACPoint, totalDrinksConsumed int, timeStepMins int) models.BACSummary {
	if len(timeline) == 0 {
		return models.BACSummary{}
	}

	var maxBAC float64
	var maxBACTime time.Time
	var soberSinceTime time.Time
	var drinkingSinceTime time.Time
	var durationOverBAC int
	var wasEverIntoxicated bool

	drinkingSinceTime = time.Time{}

	// Find max BAC and calculate time over limit
	var lastNonZeroIndex int
	for i, point := range timeline {
		if point.BAC > 0 {
			if drinkingSinceTime.IsZero() {
				drinkingSinceTime = point.Time
			}
			wasEverIntoxicated = true
			lastNonZeroIndex = i
		}

		if point.BAC > maxBAC {
			maxBAC = point.BAC
			maxBACTime = point.Time
		}

		if point.BAC > 0.08 { // Legal driving limit in most places
			// Calculate duration until next point or use the default time step
			duration := timeStepMins
			if i < len(timeline)-1 {
				duration = int(timeline[i+1].Time.Sub(point.Time).Minutes())
			}
			durationOverBAC += duration
		}

		// Track the last point where BAC > 0
		if point.BAC > 0 {
			lastNonZeroIndex = i
		}
	}

	switch {
	case !wasEverIntoxicated:
		// Person was never intoxicated by alcohol
		soberSinceTime = timeline[0].Time
	case lastNonZeroIndex == len(timeline)-1:
		// Person is still not sober at the end of timeline
		soberSinceTime = time.Time{} // Zero time indicates "not yet sober"
	default:
		// Person became sober during the timeline
		soberSinceTime = timeline[lastNonZeroIndex+1].Time
	}

	// Calculate time to sober
	var timeToSober int64
	var estimatedSoberTime time.Time
	if len(timeline) > 0 {
		lastPoint := timeline[len(timeline)-1]
		if lastPoint.BAC > 0 {
			// Calculate how many minutes until BAC reaches 0
			minutesToSober := lastPoint.BAC / metabolismRatePerMin
			// Convert to seconds for the JSON response
			timeToSober = int64(time.Duration(minutesToSober * float64(time.Minute)).Seconds())
		}
		estimatedSoberTime = lastPoint.Time.Add(time.Duration(timeToSober) * time.Second)
	}

	return models.BACSummary{
		MaxBAC:             maxBAC,
		MaxBACTime:         maxBACTime,
		SoberSinceTime:     soberSinceTime,
		TotalDrinks:        totalDrinksConsumed,
		DrinkingSinceTime:  drinkingSinceTime,
		DurationOverBAC:    durationOverBAC,
		EstimatedSoberTime: estimatedSoberTime,
	}
}

// getBACStatus returns a string description of the BAC level
func (s *Service) getBACStatus(bac float64) models.BACStatus {
	switch {
	case bac == 0:
		return models.BACStatusSober
	case bac < 0.02:
		return models.BACStatusMinimal
	case bac < 0.04:
		return models.BACStatusLight
	case bac < 0.08:
		return models.BACStatusMild
	case bac < 0.15:
		return models.BACStatusSignificant
	case bac < 0.30:
		return models.BACStatusSevere
	default:
		return models.BACStatusDangerous
	}
}
