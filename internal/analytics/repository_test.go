package analytics

import (
	"database/sql"
	"go-sober/internal/models"
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

	// Create necessary tables
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS drink_log_details (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            type TEXT NOT NULL,
            size_value REAL NOT NULL CHECK (size_value > 0),
            size_unit TEXT NOT NULL,
            abv REAL NOT NULL CHECK (abv >= 0 AND abv <= 1),
            standard_drinks REAL DEFAULT 0 CHECK (standard_drinks >= 0)
        );

        CREATE TABLE IF NOT EXISTS drink_logs (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER NOT NULL,
            drink_details_id INTEGER NOT NULL,
            logged_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (drink_details_id) REFERENCES drink_log_details(id)
        );
    `)
	if err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}

	return NewRepository(db)
}

func TestGetDrinkStats(t *testing.T) {
	repo := setupTestDB(t)
	userID := int64(1)

	// Helper function to insert test data
	insertTestDrink := func(t *testing.T, loggedAt time.Time, standardDrinks float64) {
		_, err := repo.db.Exec(`
            INSERT INTO drink_log_details (name, type, size_value, size_unit, abv, standard_drinks)
            VALUES (?, ?, ?, ?, ?, ?)`,
			"Test Beer", "Beer", 330, "ml", 0.05, standardDrinks)
		assert.NoError(t, err)

		detailsID := int64(0)
		err = repo.db.QueryRow("SELECT last_insert_rowid()").Scan(&detailsID)
		assert.NoError(t, err)

		_, err = repo.db.Exec(`
            INSERT INTO drink_logs (user_id, drink_details_id, logged_at)
            VALUES (?, ?, ?)`,
			userID, detailsID, loggedAt)
		assert.NoError(t, err)
	}

	t.Run("daily stats", func(t *testing.T) {
		now := time.Now()
		yesterday := now.AddDate(0, 0, -1)

		// Insert test data
		insertTestDrink(t, now, 1.5)
		insertTestDrink(t, now, 2.0)
		insertTestDrink(t, yesterday, 1.0)

		stats, err := repo.GetDrinkStats(userID, models.TimePeriodDaily, yesterday, now)
		assert.NoError(t, err)
		assert.Len(t, stats, 2)

		// Verify yesterday's stats
		assert.Equal(t, 1, stats[0].DrinkCount)
		assert.Equal(t, 1.0, stats[0].TotalStandardDrinks)

		// Verify today's stats
		assert.Equal(t, 2, stats[1].DrinkCount)
		assert.Equal(t, 3.5, stats[1].TotalStandardDrinks)
	})

	t.Run("weekly stats", func(t *testing.T) {
		repo = setupTestDB(t) // Reset DB
		now := time.Now()
		lastWeek := now.AddDate(0, 0, -7)

		insertTestDrink(t, now, 1.5)
		insertTestDrink(t, lastWeek, 2.0)

		stats, err := repo.GetDrinkStats(userID, models.TimePeriodWeekly, lastWeek, now)
		assert.NoError(t, err)
		assert.Len(t, stats, 2)
	})

	t.Run("no data for period", func(t *testing.T) {
		repo = setupTestDB(t)                     // Reset DB
		startDate := time.Now().AddDate(0, -1, 0) // One month ago
		endDate := time.Now()

		stats, err := repo.GetDrinkStats(userID, models.TimePeriodDaily, startDate, endDate)
		assert.NoError(t, err)
		assert.Len(t, stats, 0)
	})

	t.Run("different user's data", func(t *testing.T) {
		repo = setupTestDB(t) // Reset DB
		now := time.Now()

		// Insert drinks for two different users
		insertTestDrink(t, now, 1.5) // For userID 1

		// Insert drink for different user
		otherUserID := int64(2)
		_, err := repo.db.Exec(`
            INSERT INTO drink_log_details (name, type, size_value, size_unit, abv, standard_drinks)
            VALUES (?, ?, ?, ?, ?, ?)`,
			"Other Beer", "Beer", 330, "ml", 0.05, 2.0)
		assert.NoError(t, err)

		detailsID := int64(0)
		err = repo.db.QueryRow("SELECT last_insert_rowid()").Scan(&detailsID)
		assert.NoError(t, err)

		_, err = repo.db.Exec(`
            INSERT INTO drink_logs (user_id, drink_details_id, logged_at)
            VALUES (?, ?, ?)`,
			otherUserID, detailsID, now)
		assert.NoError(t, err)

		// Check stats for first user
		stats, err := repo.GetDrinkStats(userID, models.TimePeriodDaily, now, now)
		assert.NoError(t, err)
		assert.Len(t, stats, 1)
		assert.Equal(t, 1, stats[0].DrinkCount)
		assert.Equal(t, 1.5, stats[0].TotalStandardDrinks)
	})
}
