// Package dto
package dto

import (
	m "github.com/alehbelskidev/job_appl_track/internal/models"
)

type UpdateApplicationDto struct {
	Company *string                 `json:"company"`
	Role    *string                 `json:"role"`
	Status  *m.JobApplicationStatus `json:"status"`
}
