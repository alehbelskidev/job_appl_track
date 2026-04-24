// Package models
package models

import (
	"time"
)

type JobApplicationStatus int

const (
	Applied JobApplicationStatus = iota
	Ghosted
	Rejected
	Connected
	FailedInterview
)

type JobApplication struct {
	ID          int                  `json:"id"`
	Company     string               `json:"company"`
	Role        string               `json:"role"`
	DateApplied time.Time            `json:"date_applied"`
	DateUpdated time.Time            `json:"date_updated"`
	Status      JobApplicationStatus `json:"status"`
}
