package drinks

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"go-sober/internal/dtos"
	"go-sober/internal/models"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *Repository {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	// Create necessary tables
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS drink_templates (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			size_value REAL NOT NULL CHECK (size_value > 0),
			size_unit TEXT NOT NULL,
			abv REAL NOT NULL CHECK (abv >= 0 AND abv <= 1),
			standard_drinks REAL DEFAULT 0 CHECK (standard_drinks >= 0),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT NULL
		);

        CREATE TABLE IF NOT EXISTS drink_log_details (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			size_value REAL NOT NULL CHECK (size_value > 0),
			size_unit TEXT NOT NULL,
			abv REAL NOT NULL CHECK (abv >= 0 AND abv <= 1),
			template_id INTEGER NULL,
			hash_key TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT NULL,
			standard_drinks REAL DEFAULT 0 CHECK (standard_drinks >= 0),
			FOREIGN KEY (template_id) REFERENCES drink_templates(id)
		);

        CREATE TABLE IF NOT EXISTS drink_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			drink_details_id INTEGER NOT NULL,
			logged_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT NULL,
			FOREIGN KEY (drink_details_id) REFERENCES drink_log_details(id)
		);

    `)
	if err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}

	return NewRepository(db)
}

func TestGetDrinkLogs(t *testing.T) {
	repo := setupTestDB(t)
	userID := int64(1)

	// Insert test data
	params := dtos.CreateDrinkLogRequest{
		Name:      "Test Beer",
		Type:      "Beer",
		SizeValue: 330,
		SizeUnit:  "ml",
		ABV:       0.05,
	}

	t.Run("successful retrieval with filters", func(t *testing.T) {
		// Create multiple drink logs
		_, err := repo.CreateDrinkLog(userID, params)
		assert.NoError(t, err)

		params.Name = "Strong Beer"
		params.ABV = 0.08
		_, err = repo.CreateDrinkLog(userID, params)
		assert.NoError(t, err)

		// Test filtering
		minABV := 0.07
		filters := dtos.DrinkLogFilters{
			MinABV:    &minABV,
			DrinkType: "Beer",
			SortBy:    "abv",
			SortOrder: "desc",
		}

		logs, total, err := repo.GetDrinkLogs(userID, 1, 10, filters)
		assert.NoError(t, err)
		assert.Equal(t, 1, total)
		assert.Len(t, logs, 1)
		assert.Equal(t, "Strong Beer", logs[0].Name)
		assert.Equal(t, 0.08, logs[0].ABV)
	})

	t.Run("pagination", func(t *testing.T) {
		// Clear existing data
		repo = setupTestDB(t)

		// Insert 5 drink logs
		for i := 0; i < 5; i++ {
			params.Name = fmt.Sprintf("Beer %d", i)
			_, err := repo.CreateDrinkLog(userID, params)
			assert.NoError(t, err)
		}

		// Test first page
		logs, total, err := repo.GetDrinkLogs(userID, 1, 2, dtos.DrinkLogFilters{})
		assert.NoError(t, err)
		assert.Equal(t, 5, total)
		assert.Len(t, logs, 2)

		// Test second page
		logs, total, err = repo.GetDrinkLogs(userID, 2, 2, dtos.DrinkLogFilters{})
		assert.NoError(t, err)
		assert.Equal(t, 5, total)
		assert.Len(t, logs, 2)
	})

	t.Run("date range filtering", func(t *testing.T) {
		repo = setupTestDB(t)

		// Create drink logs with different dates
		now := time.Now()
		yesterday := now.AddDate(0, 0, -1)
		params.LoggedAt = &yesterday
		_, err := repo.CreateDrinkLog(userID, params)
		assert.NoError(t, err)

		params.LoggedAt = &now
		_, err = repo.CreateDrinkLog(userID, params)
		assert.NoError(t, err)

		// Test filtering by date range
		startDate := now.AddDate(0, 0, -1)
		endDate := now
		filters := dtos.DrinkLogFilters{
			StartDate: &startDate,
			EndDate:   &endDate,
		}

		logs, total, err := repo.GetDrinkLogs(userID, 1, 10, filters)
		assert.NoError(t, err)
		assert.Equal(t, 2, total)
		assert.Len(t, logs, 2)
	})
}

func TestUpdateDrinkLog(t *testing.T) {
	repo := setupTestDB(t)
	userID := int64(1)

	t.Run("successful update", func(t *testing.T) {
		// Create initial drink log
		params := dtos.CreateDrinkLogRequest{
			Name:      "Original Beer",
			Type:      "Beer",
			SizeValue: 330,
			SizeUnit:  "ml",
			ABV:       0.05,
		}

		logID, err := repo.CreateDrinkLog(userID, params)
		assert.NoError(t, err)

		// Update the drink log
		updateParams := dtos.UpdateDrinkLogRequest{
			ID:        logID,
			Name:      "Updated Beer",
			Type:      "Beer",
			SizeValue: 500,
			SizeUnit:  "ml",
			ABV:       0.06,
		}

		err = repo.UpdateDrinkLog(userID, updateParams)
		assert.NoError(t, err)

		// Verify the update
		logs, _, err := repo.GetDrinkLogs(userID, 1, 1, dtos.DrinkLogFilters{})
		assert.NoError(t, err)
		assert.Len(t, logs, 1)
		assert.Equal(t, "Updated Beer", logs[0].Name)
		assert.Equal(t, 500, logs[0].SizeValue)
		assert.Equal(t, 0.06, logs[0].ABV)
	})

	t.Run("unauthorized update", func(t *testing.T) {
		// Try to update with wrong user ID
		wrongUserID := int64(2)
		updateParams := dtos.UpdateDrinkLogRequest{
			ID:        1,
			Name:      "Unauthorized Update",
			Type:      "Beer",
			SizeValue: 330,
			SizeUnit:  "ml",
			ABV:       0.05,
		}

		err := repo.UpdateDrinkLog(wrongUserID, updateParams)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found or unauthorized")
	})
}

func TestDeleteDrinkLog(t *testing.T) {
	repo := setupTestDB(t)
	userID := int64(1)

	t.Run("successful deletion", func(t *testing.T) {
		// Create a drink log first
		params := dtos.CreateDrinkLogRequest{
			Name:      "Delete Test Beer",
			Type:      "Beer",
			SizeValue: 330,
			SizeUnit:  "ml",
			ABV:       0.05,
		}

		logID, err := repo.CreateDrinkLog(userID, params)
		assert.NoError(t, err)

		// Delete the drink log
		err = repo.DeleteDrinkLog(logID, userID)
		assert.NoError(t, err)

		// Verify deletion
		logs, total, err := repo.GetDrinkLogs(userID, 1, 10, dtos.DrinkLogFilters{})
		assert.NoError(t, err)
		assert.Equal(t, 0, total)
		assert.Len(t, logs, 0)
	})

	t.Run("unauthorized deletion", func(t *testing.T) {
		// Create a drink log
		params := dtos.CreateDrinkLogRequest{
			Name:      "Delete Test Beer",
			Type:      "Beer",
			SizeValue: 330,
			SizeUnit:  "ml",
			ABV:       0.05,
		}

		logID, err := repo.CreateDrinkLog(userID, params)
		assert.NoError(t, err)

		// Try to delete with wrong user ID
		wrongUserID := int64(2)
		err = repo.DeleteDrinkLog(logID, wrongUserID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found or unauthorized")
	})
}

func TestDrinkTemplates(t *testing.T) {
	repo := setupTestDB(t)

	t.Run("create and get template", func(t *testing.T) {
		template := models.DrinkTemplate{
			Name:      "Test Beer",
			Type:      "Beer",
			SizeValue: 330,
			SizeUnit:  "ml",
			ABV:       0.05,
		}

		// Create template
		err := repo.CreateDrinkTemplate(&template)
		assert.NoError(t, err)
		assert.NotZero(t, template.ID)

		// Get template by ID
		retrieved, err := repo.GetDrinkTemplate(template.ID)
		assert.NoError(t, err)
		assert.Equal(t, template.Name, retrieved.Name)
		assert.Equal(t, template.Type, retrieved.Type)
		assert.Equal(t, template.ABV, retrieved.ABV)
	})

	t.Run("get all templates", func(t *testing.T) {
		// Create multiple templates
		templates := []models.DrinkTemplate{
			{Name: "Beer 1", Type: "Beer", SizeValue: 330, SizeUnit: "ml", ABV: 0.05},
			{Name: "Beer 2", Type: "Beer", SizeValue: 500, SizeUnit: "ml", ABV: 0.045},
		}

		for _, tmpl := range templates {
			err := repo.CreateDrinkTemplate(&tmpl)
			assert.NoError(t, err)
		}

		// Get all templates
		retrieved, err := repo.GetDrinkTemplates()
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(retrieved), 2)
	})

	t.Run("update template", func(t *testing.T) {
		template := models.DrinkTemplate{
			Name:      "Original Beer",
			Type:      "Beer",
			SizeValue: 330,
			SizeUnit:  "ml",
			ABV:       0.05,
		}

		// Create template
		err := repo.CreateDrinkTemplate(&template)
		assert.NoError(t, err)

		// Update template
		template.Name = "Updated Beer"
		template.ABV = 0.06
		err = repo.UpdateDrinkTemplate(template.ID, &template)
		assert.NoError(t, err)

		// Verify update
		retrieved, err := repo.GetDrinkTemplate(template.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Beer", retrieved.Name)
		assert.Equal(t, 0.06, retrieved.ABV)
	})

	t.Run("delete template", func(t *testing.T) {
		template := models.DrinkTemplate{
			Name:      "To Delete",
			Type:      "Beer",
			SizeValue: 330,
			SizeUnit:  "ml",
			ABV:       0.05,
		}

		// Create template
		err := repo.CreateDrinkTemplate(&template)
		assert.NoError(t, err)

		// Delete template
		err = repo.DeleteDrinkTemplate(template.ID)
		assert.NoError(t, err)

		// Verify deletion
		_, err = repo.GetDrinkTemplate(template.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestDrinkLogUpdatedAt(t *testing.T) {
	t.Run("updated_at is null on creation", func(t *testing.T) {
		repo := setupTestDB(t)
		// Create initial drink log
		createReq := dtos.CreateDrinkLogRequest{
			Name:      "Beer",
			Type:      "Lager",
			SizeValue: 330,
			SizeUnit:  "ml",
			ABV:       0.05,
			LoggedAt:  nil, // Use current time
		}

		userID := int64(1)
		logID, err := repo.CreateDrinkLog(userID, createReq)
		assert.NoError(t, err)
		assert.Greater(t, logID, int64(0))

		// Fetch the created log
		logs, _, err := repo.GetDrinkLogs(userID, 1, 10, dtos.DrinkLogFilters{})
		assert.NoError(t, err)
		assert.Len(t, logs, 1)
		assert.Nil(t, logs[0].UpdatedAt)
	})

	t.Run("updated_at is set when updating", func(t *testing.T) {
		repo := setupTestDB(t)
		// Create initial drink log
		createReq := dtos.CreateDrinkLogRequest{
			Name:      "Beer",
			Type:      "Lager",
			SizeValue: 330,
			SizeUnit:  "ml",
			ABV:       0.05,
		}

		userID := int64(1)
		logID, err := repo.CreateDrinkLog(userID, createReq)
		assert.NoError(t, err)

		// Update the drink log
		updateTime := time.Now().UTC()
		updateReq := dtos.UpdateDrinkLogRequest{
			ID:        logID,
			Name:      "Craft Beer",
			Type:      "IPA",
			SizeValue: 330,
			SizeUnit:  "ml",
			ABV:       0.065,
			UpdatedAt: &updateTime,
		}

		err = repo.UpdateDrinkLog(userID, updateReq)
		assert.NoError(t, err)

		// Fetch the updated log
		logs, _, err := repo.GetDrinkLogs(userID, 1, 10, dtos.DrinkLogFilters{})
		assert.NoError(t, err)
		assert.Len(t, logs, 1)
		assert.NotNil(t, logs[0].UpdatedAt)
		assert.Equal(t, updateTime.Unix(), logs[0].UpdatedAt.Unix())
	})

	t.Run("updated_at is set even if not provided in update", func(t *testing.T) {
		repo := setupTestDB(t)
		// Create initial drink log
		createReq := dtos.CreateDrinkLogRequest{
			Name:      "Beer Miam Miam",
			Type:      "Lager",
			SizeValue: 330,
			SizeUnit:  "ml",
			ABV:       0.05,
		}

		userID := int64(2)
		logID, err := repo.CreateDrinkLog(userID, createReq)
		assert.NoError(t, err)

		// Update the drink log without updated_at
		updateReq := dtos.UpdateDrinkLogRequest{
			ID:        logID,
			Name:      "Craft Beer",
			Type:      "IPA",
			SizeValue: 330,
			SizeUnit:  "ml",
			ABV:       0.065,
		}

		err = repo.UpdateDrinkLog(userID, updateReq)
		assert.NoError(t, err)

		// Fetch the updated log
		logs, _, err := repo.GetDrinkLogs(userID, 1, 10, dtos.DrinkLogFilters{})
		assert.NoError(t, err)
		assert.Len(t, logs, 1)
		assert.NotNil(t, logs[0].UpdatedAt)
	})
}
