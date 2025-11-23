package models

type Team struct {
	TeamName string
	Members  []User
}

type TeamStats struct {
	TeamName           string
	PullRequests       int
	OpenPullRequests   int
	MergedPullRequests int
	Users              int
	ActiveUsers        int
	InactiveUsers      int
}
