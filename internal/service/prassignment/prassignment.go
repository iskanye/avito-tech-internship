package prassignment

import (
	"context"
	"log/slog"
	"time"

	"github.com/iskanye/avito-tech-internship/internal/models"
)

type PRAssignment struct {
	log *slog.Logger

	// Менеджер транзакций
	txManager TransactionManager

	// Объекты для взаимодействия с пользователями
	userCreator  UserCreator
	userModifier UserModifier
	userProvider UserProvider

	// Объекты для взаимодействия с командами
	teamCreator    TeamCreator
	teamProvider   TeamProvider
	teamModifier   TeamModifier
	teamStatistics TeamStatistics

	// Объекты для взаимодействия с пул реквестами
	prCreator  PRCreator
	prModifier PRModifier
	prProvider PRProvider

	// Объекты для взаимодействия с ревьюверами
	revAssigner ReviewersAssigner
	revModifier ReviewersModifier
}

// Менеджер транзакций
type TransactionManager interface {
	Do(
		ctx context.Context,
		f func(context.Context) error,
	) error
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

type TeamModifier interface {
	DeactivateTeam(
		ctx context.Context,
		teamName string,
	) error
}

type TeamStatistics interface {
	GetTeamsPullRequests(
		ctx context.Context,
		teamName string,
	) (
		pullRequests int,
		openPullRequests int,
		mergedPullRequests int,
		err error,
	)
}

type PRCreator interface {
	CreatePullRequest(
		ctx context.Context,
		pullRequest models.PullRequest,
	) error
}

type PRProvider interface {
	GetPullRequest(
		ctx context.Context,
		pullRequestID string,
	) (models.PullRequest, error)
	GetReview(
		ctx context.Context,
		userID string,
	) ([]models.PullRequest, error)
}

type PRModifier interface {
	MergePullRequest(
		ctx context.Context,
		pullRequestID string,
		mergedAt time.Time,
	) error
}

type ReviewersAssigner interface {
	AssignReviewers(
		ctx context.Context,
		pullRequestID string,
		authorID string,
	) error
}

type ReviewersModifier interface {
	ReassignReviewer(
		ctx context.Context,
		pullRequestID string,
		oldReviewerID string,
	) (string, error)
}

func New(
	log *slog.Logger,
	txManager TransactionManager,

	userCreator UserCreator,
	userModifier UserModifier,
	userProvider UserProvider,

	teamCreator TeamCreator,
	teamProvider TeamProvider,
	teamModifier TeamModifier,
	teamStatistics TeamStatistics,

	prCreator PRCreator,
	prModifier PRModifier,
	prProvider PRProvider,

	revAssigner ReviewersAssigner,
	revModifier ReviewersModifier,
) *PRAssignment {
	return &PRAssignment{
		log:       log,
		txManager: txManager,

		userCreator:  userCreator,
		userModifier: userModifier,
		userProvider: userProvider,

		teamCreator:    teamCreator,
		teamProvider:   teamProvider,
		teamModifier:   teamModifier,
		teamStatistics: teamStatistics,

		prCreator:  prCreator,
		prModifier: prModifier,
		prProvider: prProvider,

		revAssigner: revAssigner,
		revModifier: revModifier,
	}
}
