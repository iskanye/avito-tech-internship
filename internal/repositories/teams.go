package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/iskanye/avito-tech-internship/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Код ошибки, если в БД вносятся повторяющиеся значения
// для которых прописано UNIQUE
const UNIQUE_VIOLATION_CODE = "23505"

// Вносит имя команды в БД и возвращает её ID
func (s *Storage) AddTeam(
	ctx context.Context,
	teamName string,
) (int64, error) {
	const op = "repositories.postgres.AddTeam"

	conn := s.getter.DefaultTrOrDB(ctx, s.pool)

	// Вставить команду в базу
	insertID := conn.QueryRow(
		ctx,
		"INSERT INTO teams (team_name) VALUES ($1) RETURNING id;",
		teamName,
	)

	var id int64
	err := insertID.Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == UNIQUE_VIOLATION_CODE {
			return 0, fmt.Errorf("%s: %w", op, ErrTeamExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// Получает команду по ее названию
func (s *Storage) GetTeam(
	ctx context.Context,
	teamName string,
) (models.Team, error) {
	const op = "repositories.postgres.GetTeam"

	conn := s.getter.DefaultTrOrDB(ctx, s.pool)

	// Получаем ID команды по её названию
	getTeamID := conn.QueryRow(
		ctx,
		"SELECT id FROM teams WHERE team_name = $1;",
		teamName,
	)

	var teamID int64
	err := getTeamID.Scan(&teamID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Team{}, fmt.Errorf("%s: %w", op, ErrNotFound)
		}
		return models.Team{}, fmt.Errorf("%s: %w", op, err)
	}

	// Получаем членов команды
	getTeamMemdbers, err := conn.Query(
		ctx,
		`
		SELECT i.user_id, u.username, u.is_active
		FROM users u
		JOIN users_id i ON u.user_id = i.id
		WHERE u.team_id = $1;
		`,
		teamID,
	)
	if err != nil {
		return models.Team{}, fmt.Errorf("%s: %w", op, err)
	}
	defer getTeamMemdbers.Close()

	team := models.Team{
		TeamName: teamName,
	}
	members := make([]models.User, 0, 8)
	for getTeamMemdbers.Next() {
		var member models.User

		err = getTeamMemdbers.Scan(&member.UserID, &member.Username, &member.IsActive)
		if err != nil {
			return models.Team{}, fmt.Errorf("%s: %w", op, err)
		}

		members = append(members, member)
	}

	team.Members = members
	return team, nil
}

// Деактивирует пользователей команды
func (s *Storage) DeactivateTeam(
	ctx context.Context,
	teamName string,
) error {
	const op = "repositories.postgres.DeactivateTeam"

	conn := s.getter.DefaultTrOrDB(ctx, s.pool)

	// Получаем ID команды
	getTeamID := conn.QueryRow(
		ctx,
		`
		SELECT id
		FROM teams
		WHERE team_name = $1
		`,
		teamName,
	)

	var teamID int64
	err := getTeamID.Scan(&teamID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%s: %w", op, ErrNotFound)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	// Деактивируем всех пользователей данной команды
	_, err = conn.Exec(
		ctx,
		`
		UPDATE users
		SET is_active = FALSE
		WHERE team_id = $1;
		`,
		teamID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Получает статистику пул реквестов команды
func (s *Storage) GetTeamsPullRequests(
	ctx context.Context,
	teamName string,
) (int, int, int, error) {
	const op = "repositories.postgres.GetTeamsPullRequests"

	conn := s.getter.DefaultTrOrDB(ctx, s.pool)

	// Получаем статусы пул реквестов команды
	getPRStatuses, err := conn.Query(
		ctx,
		`
		SELECT p.status
		FROM pull_requests p
		JOIN users u ON p.author_id = u.id
		JOIN teams t ON u.team_id = t.id
		WHERE t.team_name = $1
		`,
		teamName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, 0, 0, fmt.Errorf("%s: %w", op, ErrNotFound)
		}

		return 0, 0, 0, fmt.Errorf("%s: %w", op, err)
	}
	defer getPRStatuses.Close()

	// Считаем их количество
	var pullRequests, openPullRequests int
	for getPRStatuses.Next() {
		var status string
		err := getPRStatuses.Scan(&status)
		if err != nil {
			return 0, 0, 0, fmt.Errorf("%s: %w", op, err)
		}

		pullRequests++
		if status == models.PULLREQUEST_OPEN {
			openPullRequests++
		}
	}

	mergedPullRequests := pullRequests - openPullRequests
	return pullRequests, openPullRequests, mergedPullRequests, nil
}
