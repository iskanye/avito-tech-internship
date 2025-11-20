package prassignment

import (
	"context"

	"github.com/iskanye/avito-tech-internship/internal/models"
)

type PRAssignment struct {
	// Объекты для взаимодействия с пользователями
	userCreator  UserCreator
	userModifier UserModifier
	userProvider UserProvider

	// Объекты для взаимодействия с командами
	teamCreator  TeamCreator
	teamProvider TeamProvider
}

// Интерфейсы для работы сервиса

type UserCreator interface {
	AddUser(
		ctx context.Context,
		user models.User,
	) error
}

type UserProvider interface {
	GetUser(
		ctx context.Context,
		userID string,
	) (models.User, error)
}

type UserModifier interface {
	SetActive(
		ctx context.Context,
		userID string,
		isActive bool,
	) error
}

type TeamCreator interface {
	AddTeam(
		ctx context.Context,
		teamName string,
	) (int64, error)
}

type TeamProvider interface {
	GetTeam(
		ctx context.Context,
		teamName string,
	) (models.Team, error)
}

func New()
