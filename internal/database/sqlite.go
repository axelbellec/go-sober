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
	slog.Debug("Database schema")
	for _, schema := range schemas {
		slog.Debug(schema)
	}

	return db, nil
}

// TableInfo holds information about a database table
type TableInfo struct {
	Name     string
	RowCount int
}

func ListDBSchema(db *sql.DB) ([]string, error) {
	// Query to get all table names first
	query := `
		SELECT name
		FROM sqlite_master 
		WHERE type='table' AND name NOT LIKE 'sqlite_%';
	`

	rows, err := db.Query(query)
	if err != nil {
		slog.Error("Failed to query database tables", "error", err)
		return nil, err
	}
	defer rows.Close()

	var tables []TableInfo
	for rows.Next() {
		var info TableInfo
		if err := rows.Scan(&info.Name); err != nil {
			slog.Error("Failed to scan table info", "error", err)
			return nil, err
		}

		// Get row count for each table using a separate query
		countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %q", info.Name)
		if err := db.QueryRow(countQuery).Scan(&info.RowCount); err != nil {
			slog.Error("Failed to get row count", "table", info.Name, "error", err)
			return nil, err
		}

		tables = append(tables, info)
	}

	if err = rows.Err(); err != nil {
		slog.Error("Error iterating over table rows", "error", err)
		return nil, err
	}

	// Convert TableInfo slice to string slice for backward compatibility
	results := make([]string, len(tables))
	for i, table := range tables {
		results[i] = fmt.Sprintf("%s: %d rows", table.Name, table.RowCount)
	}

	return results, nil
}
