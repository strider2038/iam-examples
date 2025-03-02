package di

import (
	"database/sql"

	"github.com/muonsoft/errors"
)

func initDB(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS products (
		id TEXT PRIMARY KEY, 
		company_id TEXT, 
		name TEXT, 
		created_by TEXT, 
		created_at DATETIME, 
		updated_at DATETIME)
	`)
	if err != nil {
		return errors.Errorf("create products table: %w", err)
	}

	return nil
}
