package auth

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *Repository {
	// Create an in-memory SQLite database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	// Create users table
	_, err = db.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create users table: %v", err)
	}

	return NewRepository(db)
}

func TestCreateUser(t *testing.T) {
	repo := setupTestDB(t)

	t.Run("successful creation", func(t *testing.T) {
		email := "test@example.com"
		password := "password123"

		err := repo.CreateUser(email, password)
		assert.NoError(t, err)

		// Verify user was created
		user, err := repo.GetUserByEmail(email)
		assert.NoError(t, err)
		assert.Equal(t, email, user.Email)
	})

	t.Run("duplicate email", func(t *testing.T) {
		email := "duplicate@example.com"
		password := "password123"

		err := repo.CreateUser(email, password)
		assert.NoError(t, err)

		// Try to create user with same email
		err = repo.CreateUser(email, password)
		assert.Error(t, err) // Should fail due to unique constraint
	})
}

func TestGetUserByEmail(t *testing.T) {
	repo := setupTestDB(t)

	t.Run("existing user", func(t *testing.T) {
		email := "test@example.com"
		password := "password123"

		// Create user first
		err := repo.CreateUser(email, password)
		assert.NoError(t, err)

		// Test retrieval
		user, err := repo.GetUserByEmail(email)
		assert.NoError(t, err)
		assert.Equal(t, email, user.Email)
		assert.NotEmpty(t, user.Password)
		assert.NotZero(t, user.CreatedAt)
	})

	t.Run("non-existent user", func(t *testing.T) {
		_, err := repo.GetUserByEmail("nonexistent@example.com")
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
	})
}

func TestComparePassword(t *testing.T) {
	repo := setupTestDB(t)

	t.Run("correct password", func(t *testing.T) {
		email := "test@example.com"
		password := "password123"

		// Create user first
		err := repo.CreateUser(email, password)
		assert.NoError(t, err)

		// Get the user to get the hashed password
		user, err := repo.GetUserByEmail(email)
		assert.NoError(t, err)

		// Test password comparison
		err = repo.ComparePassword(user.Password, password)
		assert.NoError(t, err)
	})

	t.Run("incorrect password", func(t *testing.T) {
		email := "test2@example.com"
		password := "password123"

		// Create user first
		err := repo.CreateUser(email, password)
		assert.NoError(t, err)

		// Get the user to get the hashed password
		user, err := repo.GetUserByEmail(email)
		assert.NoError(t, err)

		// Test wrong password
		err = repo.ComparePassword(user.Password, "wrongpassword")
		assert.Error(t, err)
	})
}
