package drinks

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *Repository {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	// Create drink_options table
	_, err = db.Exec(`
		CREATE TABLE drink_options (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			size_value INTEGER NOT NULL,
			size_unit TEXT NOT NULL,
			abv REAL NOT NULL
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create drink_options table: %v", err)
	}

	// Create drink_logs table
	_, err = db.Exec(`
		CREATE TABLE drink_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			drink_option_id INTEGER NOT NULL,
			logged_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (drink_option_id) REFERENCES drink_options(id)
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create drink_logs table: %v", err)
	}

	// Insert some test drink options
	_, err = db.Exec(`
		INSERT INTO drink_options (name, type, size_value, size_unit, abv)
		VALUES 
			('Beer', 'beer', 330, 'ml', 5.0),
			('Wine', 'wine', 175, 'ml', 12.0)
	`)
	if err != nil {
		t.Fatalf("Failed to insert test drink options: %v", err)
	}

	return NewRepository(db)
}

func TestGetDrinkOptions(t *testing.T) {
	repo := setupTestDB(t)

	options, err := repo.GetDrinkOptions()
	assert.NoError(t, err)
	assert.Len(t, options, 2)

	// Verify first drink option
	assert.Equal(t, "Beer", options[0].Name)
	assert.Equal(t, "beer", options[0].Type)
	assert.Equal(t, 330, options[0].SizeValue)
	assert.Equal(t, "ml", options[0].SizeUnit)
	assert.Equal(t, 5.0, options[0].ABV)
}

func TestGetDrinkOption(t *testing.T) {
	repo := setupTestDB(t)

	t.Run("existing drink option", func(t *testing.T) {
		option, err := repo.GetDrinkOption("1")
		assert.NoError(t, err)
		assert.Equal(t, "Beer", option.Name)
	})

	t.Run("non-existent drink option", func(t *testing.T) {
		_, err := repo.GetDrinkOption("999")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "drink option not found")
	})
}

func TestCreateDrinkLog(t *testing.T) {
	repo := setupTestDB(t)

	t.Run("create with current time", func(t *testing.T) {
		userID := int64(1)
		drinkOptionID := int64(1)

		id, err := repo.CreateDrinkLog(userID, drinkOptionID, nil)
		assert.NoError(t, err)
		assert.Greater(t, id, int64(0))
	})

	t.Run("create with specific time", func(t *testing.T) {
		userID := int64(1)
		drinkOptionID := int64(1)
		specificTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

		id, err := repo.CreateDrinkLog(userID, drinkOptionID, &specificTime)
		assert.NoError(t, err)
		assert.Greater(t, id, int64(0))
	})
}

func TestGetDrinkLogs(t *testing.T) {
	repo := setupTestDB(t)
	userID := int64(1)

	// Create some test logs
	specificTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	_, err := repo.CreateDrinkLog(userID, 1, &specificTime)
	assert.NoError(t, err)

	laterTime := specificTime.Add(time.Hour)
	_, err = repo.CreateDrinkLog(userID, 2, &laterTime)
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
