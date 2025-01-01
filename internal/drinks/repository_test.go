package drinks

import (
	"database/sql"
	"testing"
	"time"

	"go-sober/internal/dtos"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *Repository {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	// Create drink_templates table
	_, err = db.Exec(`
		CREATE TABLE drink_templates (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			size_value INTEGER NOT NULL,
			size_unit TEXT NOT NULL,
			abv REAL NOT NULL
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create drink_templates table: %v", err)
	}

	// Create drink_log_details table
	_, err = db.Exec(`
		CREATE TABLE drink_log_details (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			size_value INTEGER NOT NULL,
			size_unit TEXT NOT NULL,
			abv REAL NOT NULL,
			hash_key TEXT NOT NULL UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create drink_log_details table: %v", err)
	}

	// Create drink_logs table
	_, err = db.Exec(`
		CREATE TABLE drink_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			drink_details_id INTEGER NOT NULL,
			logged_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (drink_details_id) REFERENCES drink_log_details(id)
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create drink_logs table: %v", err)
	}

	// Insert some test drink templates
	_, err = db.Exec(`
		INSERT INTO drink_templates (name, type, size_value, size_unit, abv)
		VALUES 
			('Beer', 'beer', 330, 'ml', 5.0),
			('Wine', 'wine', 175, 'ml', 12.0)
	`)
	if err != nil {
		t.Fatalf("Failed to insert test drink templates: %v", err)
	}

	return NewRepository(db)
}

func TestGetDrinkTemplates(t *testing.T) {
	repo := setupTestDB(t)

	templates, err := repo.GetDrinkTemplates()
	assert.NoError(t, err)
	assert.Len(t, templates, 2)

	// Verify first drink template
	assert.Equal(t, "Beer", templates[0].Name)
	assert.Equal(t, "beer", templates[0].Type)
	assert.Equal(t, 330, templates[0].SizeValue)
	assert.Equal(t, "ml", templates[0].SizeUnit)
	assert.Equal(t, 5.0, templates[0].ABV)
}

func TestGetDrinkTemplate(t *testing.T) {
	repo := setupTestDB(t)

	t.Run("existing drink template", func(t *testing.T) {
		template, err := repo.GetDrinkTemplate(1)
		assert.NoError(t, err)
		assert.Equal(t, "Beer", template.Name)
	})

	t.Run("non-existent drink template", func(t *testing.T) {
		_, err := repo.GetDrinkTemplate(999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "drink template not found")
	})
}

func TestCreateDrinkLog(t *testing.T) {
	repo := setupTestDB(t)

	t.Run("create with current time", func(t *testing.T) {
		userID := int64(1)
		params := dtos.CreateDrinkLogRequest{
			Name:      "Beer",
			Type:      "beer",
			SizeValue: 330,
			SizeUnit:  "ml",
			ABV:       5.0,
		}

		id, err := repo.CreateDrinkLog(userID, params)
		assert.NoError(t, err)
		assert.Greater(t, id, int64(0))
	})

	t.Run("create with specific time", func(t *testing.T) {
		userID := int64(1)
		specificTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		params := dtos.CreateDrinkLogRequest{
			Name:      "Wine",
			Type:      "wine",
			SizeValue: 175,
			SizeUnit:  "ml",
			ABV:       12.0,
			LoggedAt:  &specificTime,
		}

		id, err := repo.CreateDrinkLog(userID, params)
		assert.NoError(t, err)
		assert.Greater(t, id, int64(0))
	})

	t.Run("create duplicate drink details", func(t *testing.T) {
		userID := int64(1)
		params := dtos.CreateDrinkLogRequest{
			Name:      "Beer",
			Type:      "beer",
			SizeValue: 330,
			SizeUnit:  "ml",
			ABV:       5.0,
		}

		// Create first log
		id1, err := repo.CreateDrinkLog(userID, params)
		assert.NoError(t, err)
		assert.Greater(t, id1, int64(0))

		// Create second log with same details
		id2, err := repo.CreateDrinkLog(userID, params)
		assert.NoError(t, err)
		assert.Greater(t, id2, int64(0))
		assert.NotEqual(t, id1, id2)
	})
}

func TestGetDrinkLogs(t *testing.T) {
	repo := setupTestDB(t)
	userID := int64(1)

	// Create some test logs
	specificTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	params1 := dtos.CreateDrinkLogRequest{
		Name:      "Beer",
		Type:      "beer",
		SizeValue: 330,
		SizeUnit:  "ml",
		ABV:       5.0,
		LoggedAt:  &specificTime,
	}
	_, err := repo.CreateDrinkLog(userID, params1)
	assert.NoError(t, err)

	laterTime := specificTime.Add(time.Hour)
	params2 := dtos.CreateDrinkLogRequest{
		Name:      "Wine",
		Type:      "wine",
		SizeValue: 175,
		SizeUnit:  "ml",
		ABV:       12.0,
		LoggedAt:  &laterTime,
	}
	_, err = repo.CreateDrinkLog(userID, params2)
	assert.NoError(t, err)

	t.Run("get all logs", func(t *testing.T) {
		logs, err := repo.GetDrinkLogs(userID)
		assert.NoError(t, err)
		assert.Len(t, logs, 2)
		assert.Equal(t, "Beer", logs[0].DrinkName)
		assert.Equal(t, "Wine", logs[1].DrinkName)
	})

	t.Run("get logs between dates", func(t *testing.T) {
		startTime := specificTime.Add(-time.Minute)
		endTime := specificTime.Add(30 * time.Minute)

		logs, err := repo.GetDrinkLogsBetweenDates(userID, startTime, endTime)
		assert.NoError(t, err)
		assert.Len(t, logs, 1)
		assert.Equal(t, "Beer", logs[0].DrinkName)
	})

	t.Run("get logs for different user", func(t *testing.T) {
		logs, err := repo.GetDrinkLogs(999)
		assert.NoError(t, err)
		assert.Len(t, logs, 0)
	})
}
