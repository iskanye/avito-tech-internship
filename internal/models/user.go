package models

type User struct {
	UserID   string
	Username string
	TeamID   int64
	IsActive bool
}
