package mappers

import (
	"go-sober/internal/dtos"
	"go-sober/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToBACCalculationParams(t *testing.T) {
	// Arrange
	now := time.Now()
	dto := dtos.BACCalculationRequest{
		StartTime:    now,
		EndTime:      now.Add(2 * time.Hour),
		WeightKg:     70.0,
		Gender:       "male",
		TimeStepMins: 15,
	}

	// Act
	result := ToBACCalculationParams(dto)

	// Assert
	assert.Equal(t, dto.StartTime, result.StartTime)
	assert.Equal(t, dto.EndTime, result.EndTime)
	assert.Equal(t, dto.WeightKg, result.WeightKg)
	assert.Equal(t, dto.Gender, result.Gender)
	assert.Equal(t, dto.TimeStepMins, result.TimeStepMins)
}

func TestToBACCalculationResponse(t *testing.T) {
	// Arrange
	model := models.BACCalculation{
		Timeline: []models.BACPoint{
			{
				Time: time.Now(),
				BAC:  0.08,
			},
		},
		Summary: models.BACSummary{
			MaxBAC:            0.08,
			MaxBACTime:        time.Now(),
			SoberSinceTime:    time.Now(),
			TotalDrinks:       1,
			DrinkingSinceTime: time.Now(),
			DurationOverBAC:   1,
		},
	}

	// Act
	result := ToBACCalculationResponse(model)

	// Assert
	assert.Equal(t, model.Timeline, result.Timeline)
	assert.Equal(t, model.Summary, result.Summary)
}
