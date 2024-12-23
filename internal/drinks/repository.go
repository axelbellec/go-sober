// internal/drinks/repository.go
package drinks

import (
	"database/sql"
	"fmt"
	"time"

	"go-sober/internal/models"
)

var (
	MinTime = time.Time{} // zero value
	MaxTime = time.Unix(1<<63-1, 999999999)
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetDrinkOptions() ([]models.DrinkOption, error) {
	query := `
        SELECT id, name, type, size_value, size_unit, abv
        FROM drink_options
    `

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var drinkOptions []models.DrinkOption
	for rows.Next() {
		var option models.DrinkOption
		err := rows.Scan(
			&option.ID,
			&option.Name,
			&option.Type,
			&option.SizeValue,
			&option.SizeUnit,
			&option.ABV,
		)
		if err != nil {
			return nil, err
		}
		drinkOptions = append(drinkOptions, option)
	}

	return drinkOptions, nil
}

func (r *Repository) GetDrinkOption(id string) (*models.DrinkOption, error) {
	var drinkOption models.DrinkOption
	err := r.db.QueryRow("SELECT id, name, type, size_value, size_unit, abv FROM drink_options WHERE id = ?", id).
		Scan(&drinkOption.ID, &drinkOption.Name, &drinkOption.Type, &drinkOption.SizeValue, &drinkOption.SizeUnit, &drinkOption.ABV)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("drink option not found")
		}
		return nil, err
	}
	return &drinkOption, nil
}

// Add after GetDrinkOption method

func (r *Repository) UpdateDrinkOption(id int, option *models.DrinkOption) error {
	query := `
        UPDATE drink_options 
        SET name = ?, type = ?, size_value = ?, size_unit = ?, abv = ?
        WHERE id = ?
    `

	result, err := r.db.Exec(query,
		option.Name,
		option.Type,
		option.SizeValue,
		option.SizeUnit,
		option.ABV,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to update drink option: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("drink option not found")
	}

	return nil
}

func (r *Repository) DeleteDrinkOption(id int) error {
	result, err := r.db.Exec("DELETE FROM drink_options WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete drink option: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("drink option not found")
	}

	return nil
}

func (r *Repository) CreateDrinkLog(userID int64, drinkOptionID int64, loggedAt *time.Time) (int64, error) {
	var timestamp time.Time
	if loggedAt == nil {
		timestamp = time.Now().UTC()
	} else {
		timestamp = loggedAt.UTC()
	}

	query := `
        INSERT INTO drink_logs (user_id, drink_option_id, logged_at)
        VALUES (?, ?, ?)
    `

	result, err := r.db.Exec(query, userID, drinkOptionID, timestamp)
	if err != nil {
		return 0, fmt.Errorf("failed to create drink log: %w", err)
	}

	// Get the last inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return id, nil
}

func (r *Repository) GetDrinkLogs(userID int64) ([]models.DrinkLog, error) {
	return r.GetDrinkLogsBetweenDates(userID, MinTime, MaxTime)
}

func (r *Repository) GetDrinkLogsBetweenDates(userID int64, startTime, endTime time.Time) ([]models.DrinkLog, error) {
	query := `
		SELECT 
			dl.id as id, 
			dl.user_id as user_id, 
			dl.drink_option_id as drink_option_id, 
			dl.logged_at as logged_at,
			do.name as drink_name, 
			do.abv as abv,
			do.size_value as size_value,
			do.size_unit as size_unit
		FROM drink_logs dl
		JOIN drink_options do ON dl.drink_option_id = do.id
		WHERE dl.user_id = ? 
		AND dl.logged_at >= ? 
		AND dl.logged_at < ?
		ORDER BY dl.logged_at ASC`

	utcStartTime := startTime.UTC().Format(time.DateTime)
	utcEndTime := endTime.UTC().Format(time.DateTime)

	rows, err := r.db.Query(query, userID, utcStartTime, utcEndTime)
	if err != nil {
		return nil, fmt.Errorf("error querying drink logs: %w", err)
	}
	defer rows.Close()

	var drinkLogs []models.DrinkLog
	for rows.Next() {
		var log models.DrinkLog
		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.DrinkOptionID,
			&log.LoggedAt,
			&log.DrinkName,
			&log.ABV,
			&log.SizeValue,
			&log.SizeUnit,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning drink log: %w", err)
		}
		drinkLogs = append(drinkLogs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating drink logs: %w", err)
	}

	return drinkLogs, nil
}
