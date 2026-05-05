package applications

type JobApplicationStatus int

const (
	Applied JobApplicationStatus = iota
	Ghosted
	Rejected
	Connected
	FailedInterview
)

type CreateJobApplicationDto struct {
	Company     string `json:"company"`
	Role        string `json:"role"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Notes       string `json:"notes"`
}
