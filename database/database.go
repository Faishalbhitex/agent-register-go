package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite3", "./agents.db")
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	// Create agents table if not exists
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS agents (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		skills TEXT NOT NULL,
		description TEXT,
		url TEXT NOT NULL UNIQUE,
		status TEXT DEFAULT 'registered',
		last_seen_at DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err = DB.Exec(createTableSQL); err != nil {
		return err
	}

	// Migrate existing data - add new columns if they don't exist
	migrationSQL := []string{
		`ALTER TABLE agents ADD COLUMN status TEXT DEFAULT 'registered';`,
		`ALTER TABLE agents ADD COLUMN last_seen_at DATETIME;`,
	}

	for _, sql := range migrationSQL {
		// Ignore errors for already existing columns
		DB.Exec(sql)
	}

	log.Println("Database initialized successfully")
	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
