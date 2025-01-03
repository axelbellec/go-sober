package models

// TimePeriod represents the supported time aggregation periods
type TimePeriod string

const (
	TimePeriodDaily   TimePeriod = "daily"
	TimePeriodWeekly  TimePeriod = "weekly"
	TimePeriodMonthly TimePeriod = "monthly"
	TimePeriodYearly  TimePeriod = "yearly"
	TimePeriodUnknown TimePeriod = "unknown"
)

type TimePeriodDateFormatter string

const (
	DailyDateFormatter   TimePeriodDateFormatter = "%Y-%m-%d" // 2025-07-25
	WeeklyDateFormatter  TimePeriodDateFormatter = "%Y-%W"    // 2025-30
	MonthlyDateFormatter TimePeriodDateFormatter = "%Y-%m"    // 2025-07
	YearlyDateFormatter  TimePeriodDateFormatter = "%Y"       // 2025
)

// DrinkStatsPoint represents a single data point in the statistics
type DrinkStatsPoint struct {
	TimePeriod          string  `json:"time_period"`           // The time period this stat represents
	DrinkCount          int     `json:"drink_count"`           // Number of drinks in this period
	TotalStandardDrinks float64 `json:"total_standard_drinks"` // Total standard drinks in this period
}

func ToTimePeriod(period string) TimePeriod {
	switch period {
	case "daily":
		return TimePeriodDaily
	case "weekly":
		return TimePeriodWeekly
	case "monthly":
		return TimePeriodMonthly
	case "yearly":
		return TimePeriodYearly
	}
	return TimePeriodUnknown
}

func ToTimePeriodDateFormatter(period TimePeriod) TimePeriodDateFormatter {
	switch period {
	case TimePeriodDaily:
		return DailyDateFormatter
	case TimePeriodWeekly:
		return WeeklyDateFormatter
	case TimePeriodMonthly:
		return MonthlyDateFormatter
	case TimePeriodYearly:
		return YearlyDateFormatter
	}
	return DailyDateFormatter
}
