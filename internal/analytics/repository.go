package analytics

import (
	"database/sql"
	"fmt"
	"go-sober/internal/dtos"
	"go-sober/internal/models"
	"sort"
	"strconv"
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

func (r *Repository) GetMonthlyBACStats(userID int64, startDate, endDate time.Time) ([]dtos.MonthlyBACStats, error) {
	query := `
        WITH RECURSIVE 
        months AS (
            SELECT date(?) as month
            UNION ALL
            SELECT date(month, '+1 month')
            FROM months
            WHERE month < date(?)
        ),
        numbers(n) AS (
            SELECT 0 UNION ALL SELECT n+1 FROM numbers WHERE n < 31
        ),
        days AS (
            SELECT date(month, '+' || numbers.n || ' days') as day
            FROM months
            CROSS JOIN numbers
            WHERE 
                -- Only include dates that are within the same month
                strftime('%m', date(month, '+' || numbers.n || ' days')) = strftime('%m', month)
                -- And ensure the date is valid (will exclude invalid dates like Feb 30)
                AND date(month, '+' || numbers.n || ' days') IS NOT NULL
                -- Exclude future dates
                AND date(month, '+' || numbers.n || ' days') <= date('now')
        )
        SELECT 
            strftime('%Y', d.day) as year,
            strftime('%m', d.day) as month,
            COALESCE(dl.bac_category, 'sober') as category,
            COUNT(*) as count
        FROM days d
        LEFT JOIN (
            SELECT 
                date(dl.logged_at) as log_date,
                CASE 
                    WHEN SUM(dld.standard_drinks) = 0 THEN 'sober'
                    WHEN SUM(dld.standard_drinks) < 4 THEN 'light'
                    ELSE 'heavy'
                END as bac_category
            FROM drink_logs dl
            JOIN drink_log_details dld ON dl.drink_details_id = dld.id
            WHERE dl.user_id = ?
            GROUP BY date(dl.logged_at)
        ) dl ON d.day = dl.log_date
        GROUP BY year, month, category
        ORDER BY year, month, category`

	rows, err := r.db.Query(query, startDate, endDate, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly BAC stats: %w", err)
	}
	defer rows.Close()

	statsMap := make(map[string]dtos.MonthlyBACStats)

	for rows.Next() {
		var year, month string
		var category models.BACCategory
		var count int

		if err := rows.Scan(&year, &month, &category, &count); err != nil {
			return nil, fmt.Errorf("failed to scan monthly BAC stats: %w", err)
		}

		key := fmt.Sprintf("%s-%s", year, month)
		stats, exists := statsMap[key]
		if !exists {
			y, _ := strconv.Atoi(year)
			m, _ := strconv.Atoi(month)
			stats = dtos.MonthlyBACStats{
				Year:   y,
				Month:  m,
				Counts: make(map[models.BACCategory]int),
			}
		}

		stats.Counts[category] = count
		stats.Total += count
		statsMap[key] = stats
	}

	// Convert map to slice
	var result []dtos.MonthlyBACStats
	for _, stats := range statsMap {
		result = append(result, stats)
	}

	// Sort by year and month
	sort.Slice(result, func(i, j int) bool {
		if result[i].Year != result[j].Year {
			return result[i].Year < result[j].Year
		}
		return result[i].Month < result[j].Month
	})

	return result, nil
}
