package database

import (
	"database/sql"
	"fmt"
	"go-sober/platform"
	"log/slog"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func NewSQLiteDB(dbConfig platform.DatabaseConfig) (*sql.DB, error) {

	// Check if the database file exists
	if _, err := os.Stat(dbConfig.SQL.FilePath); os.IsNotExist(err) {
		slog.Error("Database file does not exist", "path", dbConfig.SQL.FilePath)
		return nil, fmt.Errorf("database file does not exist: %s", dbConfig.SQL.FilePath)
	} else {
		slog.Debug("Database file exists", "path", dbConfig.SQL.FilePath)
	}

	db, err := sql.Open("sqlite3", dbConfig.SQL.FilePath)
	if err != nil {
		slog.Error("Failed to open database", "error", err)
		return nil, err
	}

	// Verify database connection
	if err = db.Ping(); err != nil {
		slog.Error("Failed to connect to database", "error", err)
		return nil, err
	} else {
		slog.Debug("Connected to database")
	}

	slog.Debug("Setting connection pool settings", "maxOpenConns", dbConfig.SQL.MaxOpenConns, "maxIdleConns", dbConfig.SQL.MaxIdleConns, "connMaxLifetime", dbConfig.SQL.ConnMaxLifetime)

	// Set connection pool settings
	db.SetMaxOpenConns(dbConfig.SQL.MaxOpenConns)       // Limit max open connections
	db.SetMaxIdleConns(dbConfig.SQL.MaxIdleConns)       // Limit idle connections
	db.SetConnMaxLifetime(dbConfig.SQL.ConnMaxLifetime) // Connection max lifetime

	schemas, err := ListDBSchema(db)
	if err != nil {
		slog.Error("Failed to list database schema", "error", err)
		return nil, err
	}
	slog.Debug("Database schema", "schemas", schemas)

	return db, nil
}

// ListDBSchema returns the schema of all tables in the SQLite database
func ListDBSchema(db *sql.DB) ([]string, error) {
	// Query to get all table schemas
	query := `
		SELECT sql 
		FROM sqlite_master 
		WHERE type='table' AND name NOT LIKE 'sqlite_%';
	`

	rows, err := db.Query(query)
	if err != nil {
		slog.Error("Failed to query database schema", "error", err)
		return nil, err
	}
	defer rows.Close()

	var schemas []string
	for rows.Next() {
		var schema string
		if err := rows.Scan(&schema); err != nil {
			slog.Error("Failed to scan schema row", "error", err)
			return nil, err
		}
		schemas = append(schemas, schema)
	}

	if err = rows.Err(); err != nil {
		slog.Error("Error iterating over schema rows", "error", err)
		return nil, err
	}

	return schemas, nil
}
