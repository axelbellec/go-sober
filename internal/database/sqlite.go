package database

import (
	"database/sql"
	"go-sober/platform"

	_ "github.com/mattn/go-sqlite3"
)

func NewSQLiteDB(dbConfig platform.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbConfig.SQL.FilePath)
	if err != nil {
		return nil, err
	}

	// Verify database connection
	if err = db.Ping(); err != nil {
		return nil, err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(dbConfig.SQL.MaxOpenConns)       // Limit max open connections
	db.SetMaxIdleConns(dbConfig.SQL.MaxIdleConns)       // Limit idle connections
	db.SetConnMaxLifetime(dbConfig.SQL.ConnMaxLifetime) // Connection max lifetime

	return db, nil
}
