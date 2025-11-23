package prassignment

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/iskanye/avito-tech-internship/internal/models"
	"github.com/iskanye/avito-tech-internship/internal/repositories"
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
	err := a.txManager.Do(ctx, func(ctx context.Context) error {
		// Вставляем саму команду в БД
		teamID, err := a.teamCreator.AddTeam(ctx, team.TeamName)
		if err != nil {
			log.Error("Failed to add team",
				slog.String("err", err.Error()),
			)
			if errors.Is(err, repositories.ErrTeamExists) {
				return ErrTeamExists
			}

			return fmt.Errorf("%s: %w", op, err)
		}

		// Вставляем членов команды
		for _, user := range team.Members {
			user.TeamID = teamID
			// Добавляем члена команды в БД
			err := a.userCreator.AddUser(ctx, user)
			if err != nil {
				log.Error("Failed to add members",
					slog.String("err", err.Error()),
				)

				return fmt.Errorf("%s: %w", op, err)
			}
		}

		return nil
	})
	if err != nil {
		return models.Team{}, err
	}

	log.Info("Successfully added team")

	return team, nil
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

// Получает статистику команды
func (a *PRAssignment) TeamStats(
	ctx context.Context,
	teamName string,
) (models.TeamStats, error) {
	const op = "service.PRAssignment.TeamStats"

	log := a.log.With(
		slog.String("op", op),
		slog.String("team_name", teamName),
	)

	log.Info("Attempting to get team stats")
	stats := models.TeamStats{
		TeamName: teamName,
	}

	// Получаем команду
	team, err := a.teamProvider.GetTeam(ctx, teamName)
	if err != nil {
		log.Error("Failed to get team",
			slog.String("err", err.Error()),
		)

		if errors.Is(err, repositories.ErrNotFound) {
			return models.TeamStats{}, ErrNotFound
		}
		return models.TeamStats{}, fmt.Errorf("%s: %w", op, err)
	}

	// Получаем статистику пользователей
	for _, user := range team.Members {
		stats.Users++

		if user.IsActive {
			stats.ActiveUsers++
		}
	}
	stats.InactiveUsers = stats.Users - stats.ActiveUsers

	// Получаем статистику пул реквестов
	stats.PullRequests, stats.OpenPullRequests, stats.MergedPullRequests, err = a.teamStatistics.GetTeamsPullRequests(ctx, teamName)
	if err != nil {
		log.Error("Failed to get team PR stats",
			slog.String("err", err.Error()),
		)
		return models.TeamStats{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Got team stats successfully")

	return stats, nil
}
