// Package repository
package repository

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Repo struct {
	db *sql.DB
}

func NewRepo() (*Repo, error) {
	db, err := sql.Open("sqlite3", "./local.db")
	if err != nil {
		return nil, err
	}

	return &Repo{db: db}, nil
}

func (s *Repo) InitTables() error {
	createAppsTableStmt := `
	CREATE TABLE IF NOT EXISTS job_applications (
		id INTEGER NOT NULL PRIMARY KEY,
		company TEXT NOT NULL,
		role TEXT NOT NULL,
		date_applied TEXT NOT NULL,
		date_updated TEXT,
		status INTEGER
	);
	`
	_, err := s.db.Exec(createAppsTableStmt)

	return err
}

func (s *Repo) Close() error {
	return s.db.Close()
}

func (s *Repo) CreateJobApplication(a *JobApplication) error {
	_, err := s.db.Exec(
		"INSERT INTO job_applications(company, role, date_applied, status) VALUES(?, ?, ?, ?)",
		a.Company, a.Role, a.DateApplied, a.Status,
	)
	return err
}

func (s *Repo) GetJobApplications() ([]*JobApplication, error) {
	rows, err := s.db.Query("SELECT id, company, role, date_applied, date_updated, status FROM job_applications")
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Println("Failed on deferred close of rows stmt when getting job applications")
		}
	}()

	apps := make([]*JobApplication, 0)

	for rows.Next() {
		app := &JobApplication{}
		err := rows.Scan(
			&app.ID,
			&app.Company,
			&app.Role,
			&app.DateApplied,
			&app.DateUpdated,
			&app.Status,
		)
		if err != nil {
			return nil, err
		}

		apps = append(apps, app)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return apps, nil
}
