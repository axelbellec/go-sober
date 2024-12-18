package mappers

import (
	"go-sober/internal/dtos"
	"go-sober/internal/models"
)

// ToBACCalculationParams converts BACCalculationRequest DTO to BACCalculationParams model
func ToBACCalculationParams(dto dtos.BACCalculationRequest) models.BACCalculationParams {
	return models.BACCalculationParams{
		StartTime:    dto.StartTime,
		EndTime:      dto.EndTime,
		WeightKg:     dto.WeightKg,
		Gender:       dto.Gender,
		TimeStepMins: dto.TimeStepMins,
	}
}

// ToBACCalculationResponse converts BACCalculation model to BACCalculationResponse DTO
func ToBACCalculationResponse(model models.BACCalculation) dtos.BACCalculationResponse {
	return dtos.BACCalculationResponse{
		Timeline: model.Timeline,
		Summary:  model.Summary,
	}
}
