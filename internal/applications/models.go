package applications

import (
	"time"

	"github.com/alehbelskidev/job_appl_track/internal/repo"
)

type jobApplicationStatus int

const (
	applied jobApplicationStatus = iota
	ghosted
	rejected
	connected
	failedInterview
)

type createJobApplicationDto struct {
	Company     string `json:"company"`
	Role        string `json:"role"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Notes       string `json:"notes"`
}

type importJobApplicationsDto struct {
	Company     string               `json:"company"`
	DateApplied time.Time            `json:"date_applied"`
	DateUpdated *time.Time           `json:"date_updated"`
	Description string               `json:"description"`
	Notes       string               `json:"notes"`
	Role        string               `json:"role"`
	Status      jobApplicationStatus `json:"status"`
	Url         string               `json:"url"`
}

type createApplicationResponseDto struct {
	Data repo.JobApplication `json:"data"`
}

type importApplicationsResultDto struct {
	Data []repo.JobApplication `json:"data"`
}

type createApplicationAIDto struct {
	Url   string `json:"url"`
	Notes string `json:"notes"`
}

type getJobApplicationsResponseDto struct {
	Data []repo.GetJobApplicationsRow `json:"data"`
}

type deleteApplicationResponseDto struct {
	Data bool `json:"data"`
}
