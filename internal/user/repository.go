package user

import (
	"database/sql"

	"go-sober/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) UpsertUserProfile(userID int64, profile *models.UserProfile) error {
	query := `
        INSERT INTO user_profiles (user_id, weight_kg, gender, updated_at)
        VALUES (?, ?, ?, CURRENT_TIMESTAMP)
        ON CONFLICT(user_id) DO UPDATE SET
            weight_kg = excluded.weight_kg,
            gender = excluded.gender,
            updated_at = CURRENT_TIMESTAMP
    `
	_, err := r.db.Exec(query, userID, profile.WeightKg, profile.Gender)
	return err
}

func (r *Repository) GetUserProfile(userID int64) (*models.UserProfile, error) {
	query := `
        SELECT user_id, weight_kg, gender, created_at, updated_at
        FROM user_profiles
        WHERE user_id = ?
    `
	profile := &models.UserProfile{}
	err := r.db.QueryRow(query, userID).Scan(
		&profile.UserID,
		&profile.WeightKg,
		&profile.Gender,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return profile, nil
	}
	return profile, err
}
