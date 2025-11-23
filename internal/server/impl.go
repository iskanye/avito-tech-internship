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
	DeactivateTeam(
		ctx context.Context,
		teamName string,
	) (models.Team, error)
	ReassignTeam(
		ctx context.Context,
		teamName string,
	) ([]string, error)

	// Методы пользователя
	SetIsActive(
		ctx context.Context,
		userID string,
		isActive bool,
	) (models.User, error)
	GetReview(
		ctx context.Context,
		userID string,
	) ([]models.PullRequest, error)

	// Методы пул реквестов
	CreatePullRequest(
		ctx context.Context,
		pullRequest models.PullRequest,
	) (models.PullRequest, error)
	MergePullRequest(
		ctx context.Context,
		pullRequestID string,
	) (models.PullRequest, error)
	ReassignPullRequest(
		ctx context.Context,
		pullRequestID string,
		oldReviewerId string,
	) (models.PullRequest, string, error)

	// Методы статистики
	TeamStats(
		ctx context.Context,
		teamName string,
	) (models.TeamStats, error)
}

// Проверка на реализацию всех методов
var _ api.StrictServerInterface = (*serverAPI)(nil)

func Register(engine *gin.Engine, prAssigment PRAssignment) {
	api.RegisterHandlers(engine, api.NewStrictHandler(
		&serverAPI{assign: prAssigment},
		[]api.StrictMiddlewareFunc{},
	))
}
