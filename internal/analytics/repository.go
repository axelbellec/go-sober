package analytics

import (
	"database/sql"
	"fmt"
	"go-sober/internal/models"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetDrinkStats(userID int64, period models.TimePeriod, startDate time.Time, endDate time.Time) ([]models.DrinkStatsPoint, error) {
	query := `
        SELECT 
            strftime(?, dl.logged_at) as time_period,
            COUNT(*) as drink_count,
            ROUND(SUM(dld.standard_drinks), 2) as total_standard_drinks
        FROM drink_logs dl
        JOIN drink_log_details dld ON dl.drink_details_id = dld.id
        WHERE dl.user_id = ?
        AND dl.logged_at >= ?
        AND dl.logged_at <= ?
        GROUP BY strftime(?, dl.logged_at)
        ORDER BY time_period ASC`

	periodDateFormatter := models.ToTimePeriodDateFormatter(period)

	rows, err := r.db.Query(query, periodDateFormatter, userID, startDate, endDate, periodDateFormatter)
	if err != nil {
		return nil, fmt.Errorf("failed to get drink stats: %w", err)
	}
	defer rows.Close()

	var stats []models.DrinkStatsPoint
	for rows.Next() {
		var stat models.DrinkStatsPoint
		var drinkCount int64
		var totalStandardDrinks float64

		if err := rows.Scan(&stat.TimePeriod, &drinkCount, &totalStandardDrinks); err != nil {
			return nil, fmt.Errorf("failed to scan drink stats: %w", err)
		}

		stat.DrinkCount = int(drinkCount)
		stat.TotalStandardDrinks = totalStandardDrinks
		stats = append(stats, stat)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating drink stats rows: %w", err)
	}

	return stats, nil
}
