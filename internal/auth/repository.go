package auth

import (
	"database/sql"
	"log/slog"

	"go-sober/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(email, password string) error {
	// Hash the password with bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Could not hash password", "error", err)
		return err
	}

	query := `
        INSERT INTO users (email, password)
        VALUES (?, ?)
    `
	_, err = r.db.Exec(query, email, string(hashedPassword))
	return err
}

func (r *Repository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
        SELECT id, email, password, created_at
        FROM users
        WHERE email = ?
    `
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		slog.Error("Could not get user by email", "error", err)
		return nil, err
	}
	return user, nil
}

// ComparePassword compares a hashed password with a plain text password
func (r *Repository) ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
