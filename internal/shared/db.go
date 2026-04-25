package shared

import (
	"database/sql"
	"log"
)

func InitTables(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Print(err)
		}
	}()

	createUsersTableStmt := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER NOT NULL PRIMARY KEY,
		email TEXT NOT NULL,
		password TEXT NOT NULL
	)
	`
	_, err = tx.Exec(createUsersTableStmt)
	if err != nil {
		return err
	}

	createAppsTableStmt := `
	CREATE TABLE IF NOT EXISTS job_applications (
		id INTEGER NOT NULL PRIMARY KEY,
		company TEXT NOT NULL,
		role TEXT NOT NULL,
		date_applied TEXT NOT NULL,
		date_updated TEXT,
		status INTEGER,
		owner_id INTEGER REFERENCES users(id)
	);
	`
	_, err = tx.Exec(createAppsTableStmt)
	if err != nil {
		return err
	}

	createRefreshTableStmt := `
	CREATE TABLE IF NOT EXISTS refresh_tokens (
		id INTEGER NOT NULL PRIMARY KEY,
		token TEXT NOT NULL,
	)
	`
	_, err = tx.Exec(createRefreshTableStmt)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return err
}
