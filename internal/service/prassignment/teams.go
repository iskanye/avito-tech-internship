package prassignment

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/iskanye/avito-tech-internship/internal/models"
	"github.com/iskanye/avito-tech-internship/internal/repositories"
	"golang.org/x/sync/errgroup"
)

// Создаёт команду и создаёт/обновляет её пользователей
func (a *PRAssignment) AddTeam(
	ctx context.Context,
	team models.Team,
) (models.Team, error) {
	const op = "service.PRAssignment.AddTeam"

	log := a.log.With(
		slog.String("op", op),
		slog.String("team_name", team.TeamName),
	)

	log.Info("Attempting to add team")

	// Начинаем транзакцию
	err := a.txManager.Begin(ctx)
	if err != nil {
		log.Error("Failed to start transaction",
			slog.String("err", err.Error()),
		)
		return models.Team{}, fmt.Errorf("%s: %w", op, err)
	}
	defer a.txManager.Rollback(ctx)

	// Вставляем саму команду в БД
	teamID, err := a.teamCreator.AddTeam(ctx, team.TeamName)
	if err != nil {
		log.Error("Failed to get team",
			slog.String("err", err.Error()),
		)
		if errors.Is(err, repositories.ErrTeamExists) {
			return models.Team{}, ErrTeamExists
		}

		return models.Team{}, fmt.Errorf("%s: %w", op, err)
	}

	// Параллелезируем каждую вставку в БД
	errGroup, errCtx := errgroup.WithContext(ctx)

	for _, user := range team.Members {
		user.TeamID = teamID
		errGroup.Go(func() error {
			return a.addTeamMember(errCtx, user)
		})
	}

	// Ждём выполнения всех вставок
	err = errGroup.Wait()
	if err != nil {
		log.Error("Failed to add team members",
			slog.String("err", err.Error()),
		)
		return models.Team{}, fmt.Errorf("%s: %w", op, err)
	}

	if err = a.txManager.Commit(ctx); err != nil {
		log.Error("Failed to commit transaction",
			slog.String("err", err.Error()),
		)
		return models.Team{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully added team")

	return team, nil
}

// Добавить члена команды в БД
func (a *PRAssignment) addTeamMember(
	ctx context.Context,
	member models.User,
) error {
	// Добавляем члена команды
	err := a.userCreator.AddUser(ctx, member)
	if errors.Is(err, repositories.ErrUserExists) {
		// Если член команды уже есть в БД - обновляем его данные
		return a.userModifier.UpdateUser(ctx, member)
	}

	return err
}

// Получить команду по её названию
func (a *PRAssignment) GetTeam(
	ctx context.Context,
	teamName string,
) (models.Team, error) {
	const op = "service.PRAssignment.GetTeam"

	log := a.log.With(
		slog.String("op", op),
		slog.String("team_name", teamName),
	)

	log.Info("Attempting to get team")

	// Получаем команду
	team, err := a.teamProvider.GetTeam(ctx, teamName)
	if err != nil {
		log.Error("Failed to get team",
			slog.String("err", err.Error()),
		)

		if errors.Is(err, repositories.ErrNotFound) {
			return models.Team{}, ErrNotFound
		}
		return models.Team{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully got team")

	return team, nil
}
