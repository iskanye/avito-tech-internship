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

	// Вставить команду в базу
	insertID := s.pool.QueryRow(
		ctx,
		"INSERT INTO teams (team_name) VALUES ($1) RETURNING id;",
		teamName,
	)

	var id int64
	err := insertID.Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == UNIQUE_VIOLATION_CODE {
			return 0, ErrTeamExists
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

	// Получаем ID команды по её названию
	getTeamID := s.pool.QueryRow(
		ctx,
		"SELECT id FROM teams WHERE team_name = $1;",
		teamName,
	)

	var teamID int64
	err := getTeamID.Scan(&teamID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Team{}, ErrNotFound
		}
		return models.Team{}, fmt.Errorf("%s: %w", op, err)
	}

	// Получаем членов команды
	getTeamMemdbers, err := s.pool.Query(
		ctx,
		`
		SELECT i.user_id, u.username, u.is_active
		FROM users u
		JOIN users_id i ON i.id = u.user_id
		WHERE team_id = $1;
		`,
		teamID,
	)
	if err != nil {
		return models.Team{}, fmt.Errorf("%s: %w", op, err)
	}

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
