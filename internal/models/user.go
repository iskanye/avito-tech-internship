package models

type User struct {
	UserID   string
	Username string
	TeamID   int64 // Для внесения в БД использует ID команды
	TeamName string
	IsActive bool
}
