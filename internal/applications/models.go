package applications

type jobApplicationStatus int

const (
	applied jobApplicationStatus = iota
	ghosted
	rejected
	connected
	failedInterview
)

type createJobApplicationDto struct {
	Company string `json:"company"`
	Role    string `json:"role"`
}
