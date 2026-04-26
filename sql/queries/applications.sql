-- name: GetJobApplications :many
SELECT id, company, role, date_applied, date_updated, status
FROM job_applications
WHERE owner_id=$1;

-- name: CreateJobApplication :one
INSERT INTO job_applications(
  company, role, date_applied, status, owner_id
) VALUES($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateJobApplication :one
UPDATE job_applications
SET
  company = $3,
  role = $4,
  status = $5,
  date_updated = NOW()
WHERE id=$1 AND owner_id=$2
RETURNING *;
