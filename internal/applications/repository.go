package applications

import (
	"database/sql"
	"log"
	"time"
)

type repository struct {
	db *sql.DB
}

func newRepository(db *sql.DB) *repository {
	return &repository{db: db}
}

func (r *repository) createJobApplication(a *JobApplication) error {
	_, err := r.db.Exec(
		"INSERT INTO job_applications(company, role, date_applied, status) VALUES(?, ?, ?, ?)",
		a.Company, a.Role, a.DateApplied, a.Status,
	)
	return err
}

func (r *repository) getJobApplications() ([]*JobApplication, error) {
	rows, err := r.db.Query("SELECT id, company, role, date_applied, date_updated, status FROM job_applications")
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

func (r *repository) updateApplication(id int, payload *UpdateApplicationDto) error {
	query := "UPDATE job_applications SET date_updated = ?"
	args := []any{time.Now()}

	if payload.Company != nil {
		query += ", company = ?"
		args = append(args, *payload.Company)
	}
	if payload.Role != nil {
		query += ", role = ?"
		args = append(args, *payload.Role)
	}
	if payload.Status != nil {
		query += ", status = ?"
		args = append(args, *payload.Status)
	}

	query += " WHERE id = ?"
	args = append(args, id)

	_, err := r.db.Exec(query, args...)
	return err
}
