package models

import "time"

type PRStatus = string

const (
	PULLREQUEST_OPEN   PRStatus = "OPEN"
	PULLREQUEST_MERGED PRStatus = "MERGED"
)

type PullRequest struct {
	ID                string
	Name              string
	AuthorID          string
	Status            PRStatus
	AssignedReviewers []string
	CreatedAt         time.Time
	MergedAt          time.Time
}
