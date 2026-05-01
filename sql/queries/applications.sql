-- name: GetJobApplications :many
SELECT id, company, role, date_applied, date_updated, status, description, url, notes
FROM job_applications
WHERE owner_id=$1;

-- name: CreateJobApplication :one
INSERT INTO job_applications(
  company, role, date_applied, status, owner_id, description, url, notes
) VALUES($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateJobApplication :one
UPDATE job_applications
SET
  company = $3,
  role = $4,
  status = $5,
  date_updated = NOW(),
  description = $6,
  url = $7,
  notes = $8
WHERE id=$1 AND owner_id=$2
RETURNING *;

-- name: UpdateJobApplicationStatus :one
UPDATE job_applications
SET
  status = $3,
  date_updated = NOW()
WHERE id=$1 AND owner_id=$2
RETURNING *;

-- name: DeleteJobApplication :exec
DELETE FROM job_applications
WHERE id=$1;
