-- name: GetJobApplications :many
SELECT id, company, role, date_applied, date_updated, status
FROM job_applications;

-- name: CreateJobApplication :one
INSERT INTO job_applications(
  company, role, date_applied, status
) VALUES($1, $2, $3, $4)
RETURNING *;

-- name: UpdateJobApplication :one
UPDATE job_applications
SET
  company = $2,
  role = $3,
  status = $4,
  date_updated = NOW()
WHERE id=$1
RETURNING *;
