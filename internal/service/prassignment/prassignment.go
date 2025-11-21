package prassignment

import (
	"context"
	"log/slog"

	"github.com/iskanye/avito-tech-internship/internal/models"
)

type PRAssignment struct {
	log *slog.Logger

	// Объект управления транзакциями БД
	txManager TransactionManager

	// Объекты для взаимодействия с пользователями
	userCreator  UserCreator
	userModifier UserModifier
	userProvider UserProvider

	// Объекты для взаимодействия с командами
	teamCreator  TeamCreator
	teamProvider TeamProvider

	// Объекты для взаимодействия с пул реквестами
	prCreator  PRCreator
	prModifier PRModifier
	prProvider PRProvider

	// Объекты для взаимодействия с ревьюверами
	revAssigner ReviewersAssigner
	revModifier ReviewersModifier
}

// Интерфейс транзакций
type TransactionManager interface {
	Begin(context.Context) error
	Rollback(context.Context) error
	Commit(context.Context) error
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
	UpdateUser(
		ctx context.Context,
		user models.User,
	) error
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

		teamCreator:  teamCreator,
		teamProvider: teamProvider,

		prCreator:  prCreator,
		prModifier: prModifier,
		prProvider: prProvider,

		revAssigner: revAssigner,
		revModifier: revModifier,
	}
}
