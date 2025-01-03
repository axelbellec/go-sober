package models

// TimePeriod represents the supported time aggregation periods
type TimePeriod string

const (
	Daily             TimePeriod = "daily"
	Weekly            TimePeriod = "weekly"
	Monthly           TimePeriod = "monthly"
	Yearly            TimePeriod = "yearly"
	TimePeriodUnknown TimePeriod = "unknown"
)

type TimePeriodDateFormatter string

const (
	DailyDateFormatter   TimePeriodDateFormatter = "%Y-%m-%d" // 2025-07-25
	WeeklyDateFormatter  TimePeriodDateFormatter = "%Y-%W"    // 2025-30
	MonthlyDateFormatter TimePeriodDateFormatter = "%Y-%m"    // 2025-07
	YearlyDateFormatter  TimePeriodDateFormatter = "%Y"       // 2025
)

// DrinkStats represents generic drink statistics for any time period
type DrinkStats struct {
	TimePeriod          string  `json:"time_period"`
	DrinkCount          int     `json:"drink_count"`
	TotalStandardDrinks float64 `json:"total_standard_drinks"`
}

func ToTimePeriod(period string) TimePeriod {
	switch period {
	case "daily":
		return Daily
	case "weekly":
		return Weekly
	case "monthly":
		return Monthly
	case "yearly":
		return Yearly
	}
	return TimePeriodUnknown
}

func ToTimePeriodDateFormatter(period TimePeriod) TimePeriodDateFormatter {
	switch period {
	case Daily:
		return DailyDateFormatter
	case Weekly:
		return WeeklyDateFormatter
	case Monthly:
		return MonthlyDateFormatter
	case Yearly:
		return YearlyDateFormatter
	}
	return DailyDateFormatter
}
