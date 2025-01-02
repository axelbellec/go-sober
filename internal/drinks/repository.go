// internal/drinks/repository.go
package drinks

import (
	"database/sql"
	"fmt"
	"time"

	"go-sober/internal/dtos"
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

func (r *Repository) GetDrinkTemplates() ([]models.DrinkTemplate, error) {
	query := `
        SELECT id, name, type, size_value, size_unit, abv
        FROM drink_templates
    `

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var drinkTemplates []models.DrinkTemplate
	for rows.Next() {
		var template models.DrinkTemplate
		err := rows.Scan(
			&template.ID,
			&template.Name,
			&template.Type,
			&template.SizeValue,
			&template.SizeUnit,
			&template.ABV,
		)
		if err != nil {
			return nil, err
		}
		drinkTemplates = append(drinkTemplates, template)
	}

	return drinkTemplates, nil
}

func (r *Repository) GetDrinkTemplate(id int) (*models.DrinkTemplate, error) {
	var drinkTemplate models.DrinkTemplate
	err := r.db.QueryRow("SELECT id, name, type, size_value, size_unit, abv FROM drink_templates WHERE id = ?", id).
		Scan(&drinkTemplate.ID, &drinkTemplate.Name, &drinkTemplate.Type, &drinkTemplate.SizeValue, &drinkTemplate.SizeUnit, &drinkTemplate.ABV)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("drink template not found")
		}
		return nil, err
	}
	return &drinkTemplate, nil
}

func (r *Repository) CreateDrinkTemplate(template *models.DrinkTemplate) error {
	query := `
        INSERT INTO drink_templates (name, type, size_value, size_unit, abv)
        VALUES (?, ?, ?, ?, ?)
    `

	result, err := r.db.Exec(query, template.Name, template.Type, template.SizeValue, template.SizeUnit, template.ABV)
	if err != nil {
		return fmt.Errorf("failed to create drink template: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	template.ID = int(id)
	return nil
}

func (r *Repository) UpdateDrinkTemplate(id int, template *models.DrinkTemplate) error {
	query := `
        UPDATE drink_templates 
        SET name = ?, type = ?, size_value = ?, size_unit = ?, abv = ?
        WHERE id = ?
    `

	result, err := r.db.Exec(query,
		template.Name,
		template.Type,
		template.SizeValue,
		template.SizeUnit,
		template.ABV,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to update drink template: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("drink template not found")
	}

	return nil
}

func (r *Repository) DeleteDrinkTemplate(id int) error {
	result, err := r.db.Exec("DELETE FROM drink_templates WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete drink template: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("drink template not found")
	}

	return nil
}

func (r *Repository) CreateDrinkLog(userID int64, params dtos.CreateDrinkLogRequest) (int64, error) {
	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if we return early due to an error

	// hash all the CreateDrinkLogRequest
	hashKey := fmt.Sprintf("%s-%s-%d-%s-%f", params.Name, params.Type, params.SizeValue, params.SizeUnit, params.ABV)

	// check if the hashKey is in drink_log_details
	query := `SELECT id FROM drink_log_details WHERE hash_key = ?`
	var drinkLogDetailID int64
	err = tx.QueryRow(query, hashKey).Scan(&drinkLogDetailID)
	if err != nil && err != sql.ErrNoRows {
		return 0, fmt.Errorf("failed to check if drink log detail exists: %w", err)
	}

	var loggedAt time.Time
	if params.LoggedAt == nil {
		loggedAt = time.Now().UTC()
	} else {
		loggedAt = params.LoggedAt.UTC()
	}

	// if the hashKey is not in drink_log_details, create a new drink_log_details
	if err == sql.ErrNoRows {
		query = `INSERT INTO drink_log_details (name, type, size_value, size_unit, abv, hash_key, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
		result, err := tx.Exec(query, params.Name, params.Type, params.SizeValue, params.SizeUnit, params.ABV, hashKey, loggedAt)
		if err != nil {
			return 0, fmt.Errorf("failed to create drink log detail: %w", err)
		}
		drinkLogDetailID, err = result.LastInsertId()
		if err != nil {
			return 0, fmt.Errorf("failed to get last insert id: %w", err)
		}
	}

	// create a new drink_log
	query = `INSERT INTO drink_logs (user_id, drink_details_id, logged_at) VALUES (?, ?, ?)`
	result, err := tx.Exec(query, userID, drinkLogDetailID, loggedAt)
	if err != nil {
		return 0, fmt.Errorf("failed to create drink log: %w", err)
	}
	drinkLogID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return drinkLogID, nil
}

func (r *Repository) GetDrinkLogs(userID int64, page, pageSize int, filters dtos.DrinkLogFilters) ([]models.DrinkLog, int, error) {
	// Build the base query with filters
	baseQuery := `
        SELECT 
            dl.id, dl.user_id, dl.logged_at, dl.updated_at,
            dld.name, dld.type, dld.size_value,
            dld.size_unit, dld.abv
        FROM drink_logs dl
        JOIN drink_log_details dld ON dl.drink_details_id = dld.id
        WHERE dl.user_id = ?
    `
	countQuery := `
        SELECT COUNT(*)
        FROM drink_logs dl
        JOIN drink_log_details dld ON dl.drink_details_id = dld.id
        WHERE dl.user_id = ?
    `

	args := []interface{}{userID}
	filterClauses := []string{}

	// Add date range filters
	if filters.StartDate != nil {
		filterClauses = append(filterClauses, "dl.logged_at >= ?")
		args = append(args, filters.StartDate.UTC())
	}
	if filters.EndDate != nil {
		filterClauses = append(filterClauses, "dl.logged_at <= ?")
		args = append(args, filters.EndDate.UTC())
	}

	// Add drink type filter
	if filters.DrinkType != "" {
		filterClauses = append(filterClauses, "dld.type = ?")
		args = append(args, filters.DrinkType)
	}

	// Add ABV range filters
	if filters.MinABV != nil {
		filterClauses = append(filterClauses, "dld.abv >= ?")
		args = append(args, *filters.MinABV)
	}
	if filters.MaxABV != nil {
		filterClauses = append(filterClauses, "dld.abv <= ?")
		args = append(args, *filters.MaxABV)
	}

	// Add filter clauses to queries
	for _, clause := range filterClauses {
		baseQuery += " AND " + clause
		countQuery += " AND " + clause
	}

	// Add sorting
	if filters.SortBy != "" {
		sortOrder := "ASC"
		if filters.SortOrder == "desc" {
			sortOrder = "DESC"
		}

		// Validate sort column to prevent SQL injection
		validSortColumns := map[string]string{
			"logged_at":  "dl.logged_at",
			"updated_at": "dl.updated_at",
			"abv":        "dld.abv",
			"size_value": "dld.size_value",
			"name":       "dld.name",
			"type":       "dld.type",
		}

		if sortCol, valid := validSortColumns[filters.SortBy]; valid {
			baseQuery += fmt.Sprintf(" ORDER BY %s %s", sortCol, sortOrder)
		}
	} else {
		baseQuery += " ORDER BY dl.logged_at DESC"
	}

	// Add pagination
	baseQuery += " LIMIT ? OFFSET ?"
	args = append(args, pageSize, (page-1)*pageSize)

	// Get total count first
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting drink logs: %w", err)
	}

	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying drink logs: %w", err)
	}
	defer rows.Close()

	var drinkLogs []models.DrinkLog
	for rows.Next() {
		var log models.DrinkLog
		var updatedAt sql.NullTime
		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.LoggedAt,
			&updatedAt,
			&log.Name,
			&log.Type,
			&log.SizeValue,
			&log.SizeUnit,
			&log.ABV,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning drink log: %w", err)
		}

		fmt.Printf("updatedAt: %v\n", updatedAt)
		fmt.Printf("updatedAt.Valid: %v\n", updatedAt.Valid)
		fmt.Printf("updatedAt.Time: %v\n", updatedAt.Time)

		// if updatedAt is not null, set log.UpdatedAt to updatedAt.Time
		if updatedAt.Valid {
			log.UpdatedAt = &updatedAt.Time
		}

		drinkLogs = append(drinkLogs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating drink logs: %w", err)
	}

	return drinkLogs, total, nil
}

func (r *Repository) UpdateDrinkLog(userID int64, params dtos.UpdateDrinkLogRequest) error {
	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Verify the log exists and belongs to the user
	logID := params.ID

	var oldDetailsID int64
	err = tx.QueryRow("SELECT drink_details_id FROM drink_logs WHERE id = ? AND user_id = ?", logID, userID).Scan(&oldDetailsID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("drink log not found or unauthorized")
		}
		return fmt.Errorf("failed to fetch drink log: %w", err)
	}

	// Create hash key for the new details
	hashKey := fmt.Sprintf("%s-%s-%d-%s-%f", params.Name, params.Type, params.SizeValue, params.SizeUnit, params.ABV)

	// Check if the new details already exist
	var newDetailsID int64
	err = tx.QueryRow("SELECT id FROM drink_log_details WHERE hash_key = ?", hashKey).Scan(&newDetailsID)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check drink log details: %w", err)
	}

	// If details don't exist, create them
	if err == sql.ErrNoRows {
		query := `INSERT INTO drink_log_details (name, type, size_value, size_unit, abv, hash_key, created_at) 
				 VALUES (?, ?, ?, ?, ?, ?, ?)`
		result, err := tx.Exec(query, params.Name, params.Type, params.SizeValue, params.SizeUnit, params.ABV, hashKey, time.Now().UTC())
		if err != nil {
			return fmt.Errorf("failed to create drink log details: %w", err)
		}
		newDetailsID, err = result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert id: %w", err)
		}
	}

	// Update the drink log with new details and timestamp
	updatedAt := time.Now().UTC()
	if params.UpdatedAt != nil {
		updatedAt = params.UpdatedAt.UTC()
	}

	_, err = tx.Exec("UPDATE drink_logs SET drink_details_id = ?, updated_at = ? WHERE id = ? AND user_id = ?",
		newDetailsID, updatedAt, logID, userID)
	if err != nil {
		return fmt.Errorf("failed to update drink log: %w", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *Repository) DeleteDrinkLog(logID int64, userID int64) error {
	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete the log only if it belongs to the user
	result, err := tx.Exec("DELETE FROM drink_logs WHERE id = ? AND user_id = ?", logID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete drink log: %w", err)
	}

	// Check if any row was actually deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("drink log not found or unauthorized")
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
