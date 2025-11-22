package repositories

import "errors"

var (
	ErrNotFound     = errors.New("resource not found")
	ErrTeamExists   = errors.New("team already exists")
	ErrUserExists   = errors.New("user already exists")
	ErrPRExists     = errors.New("PR already exists")
	ErrNoCandidates = errors.New("no candidate found")
)
