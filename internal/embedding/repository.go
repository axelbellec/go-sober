package embedding

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) StoreEmbedding(drinkOptionID int, embedding []float64) error {
	// Convert embedding to JSON string
	embeddingJSON, err := json.Marshal(embedding)
	if err != nil {
		return fmt.Errorf("error marshaling embedding: %w", err)
	}

	// First try to update existing record
	slog.Info("attempting to update embedding", "drink_option_id", drinkOptionID)
	query := `
        UPDATE drink_embeddings 
        SET embedding_data = ?
        WHERE drink_option_id = ?
    `
	result, err := r.db.Exec(query, string(embeddingJSON), drinkOptionID)
	if err != nil {
		return fmt.Errorf("error updating embedding: %w", err)
	}

	slog.Debug("checking rows affected")
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		slog.Info("no existing record found, inserting new embedding", "drink_option_id", drinkOptionID)
		query = `
            INSERT INTO drink_embeddings (drink_option_id, embedding_data)
            VALUES (?, ?)
        `
		_, err = r.db.Exec(query, drinkOptionID, string(embeddingJSON))
		if err != nil {
			return fmt.Errorf("error inserting embedding: %w", err)
		}
	}

	return nil
}

func (r *Repository) GetEmbedding(drinkOptionID int) ([]float64, error) {
	var embeddingJSON string
	query := `
        SELECT embedding_data 
        FROM drink_embeddings 
        WHERE drink_option_id = ?
    `

	err := r.db.QueryRow(query, drinkOptionID).Scan(&embeddingJSON)
	if err != nil {
		return nil, err
	}

	var embedding []float64
	err = json.Unmarshal([]byte(embeddingJSON), &embedding)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling embedding: %w", err)
	}

	return embedding, nil
}
