package db

import (
	"database/sql"
	"fmt"
)

type DB struct {
	*sql.DB
}

func NewDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	err = initSchema(db)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}
	return db, nil
}

func initSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    login TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS expressions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    expression TEXT NOT NULL,
    status TEXT NOT NULL,
    result REAL DEFAULT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS tasks (
    id TEXT PRIMARY KEY,
    expression_id TEXT NOT NULL,
    arg1 TEXT NOT NULL,
    arg2 TEXT NOT NULL,
    operation TEXT NOT NULL,
    status TEXT NOT NULL,
    result REAL,
    dependencies TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (expression_id) REFERENCES expressions(id)
);`

	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}
	return err
}

func CloseDB(db *DB) error {
	if err := db.DB.Close(); err != nil {
		return err
	}
	return nil
}
