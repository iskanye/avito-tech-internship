package server

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/iskanye/avito-tech-internship/internal/models"
	"github.com/iskanye/avito-tech-internship/pkg/api"
)

type serverAPI struct {
	assign PRAssignment
}

type PRAssignment interface {
	// Методы команд
	AddTeam(
		ctx context.Context,
		team models.Team,
	) (models.Team, error)
	GetTeam(
		ctx context.Context,
		teamName string,
	) (models.Team, error)
	// Методы пользователя
	SetIsActive(
		ctx context.Context,
		userID string,
		isActive bool,
	) (models.User, error)
}

// Проверка на реализацию всех методов
var _ api.StrictServerInterface = (*serverAPI)(nil)

func Register(engine *gin.Engine, prAssigment PRAssignment) {
	api.RegisterHandlers(engine, api.NewStrictHandler(
		&serverAPI{assign: prAssigment},
		[]api.StrictMiddlewareFunc{},
	))
}
