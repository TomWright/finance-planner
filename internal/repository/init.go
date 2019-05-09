package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"path/filepath"
)

// ConnectSQLite connects to an SQLite db in the given storage directory.
func ConnectSQLite(storageDir string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filepath.Join(storageDir, "finance.db"))
	if err != nil {
		return nil, fmt.Errorf("could not open sqlite db: %s", err)
	}
	return db, nil
}
