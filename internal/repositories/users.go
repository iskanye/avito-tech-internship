package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/iskanye/avito-tech-internship/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Добавляет пользователя в БД.
func (s *Storage) AddUser(
	ctx context.Context,
	user models.User,
) error {
	const op = "repositories.postgres.AddUser"

	conn := s.getter.DefaultTrOrDB(ctx, s.pool)

	// Вставляем ID в базу
	getUserID := conn.QueryRow(
		ctx,
		"INSERT INTO users_id (user_id) VALUES ($1) RETURNING id;",
		user.UserID,
	)

	var id int64
	err := getUserID.Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == UNIQUE_VIOLATION_CODE {
			return ErrUserExists
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	// Вставляем юзера
	_, err = conn.Exec(
		ctx,
		"INSERT INTO users (user_id, username, team_id, is_active) VALUES ($1, $2, $3, $4);",
		id, user.Username, user.TeamID, user.IsActive,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Обновить данные пользователя в БД
func (s *Storage) UpdateUser(
	ctx context.Context,
	user models.User,
) error {
	const op = "repositories.postgres.AddUser"

	conn := s.getter.DefaultTrOrDB(ctx, s.pool)

	// Получаем ID пользователя
	id, err := s.getUserID(ctx, user.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = conn.Exec(
		ctx,
		"UPDATE users SET username = $1, team_id = $2, is_active = $3 WHERE id = $4",
		user.Username, user.TeamID, user.IsActive, id,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Возвращает пользователя по его ID
func (s *Storage) GetUser(
	ctx context.Context,
	userID string,
) (models.User, error) {
	const op = "repositories.postgres.GetUser"

	conn := s.getter.DefaultTrOrDB(ctx, s.pool)

	// Получаем ID пользователя
	id, err := s.getUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, ErrNotFound
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	// Получаем данные о пользователе из БД
	res := conn.QueryRow(
		ctx,
		`
		SELECT u.username, t.team_name, u.is_active 
		FROM users u
		JOIN teams t ON u.team_id = t.id
		WHERE u.id = $1;
		`,
		id,
	)

	user := models.User{
		UserID: userID,
	}
	err = res.Scan(&user.Username, &user.TeamName, &user.IsActive)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// Меняет is_active у пользователя
func (s *Storage) SetActive(
	ctx context.Context,
	userID string,
	isActive bool,
) error {
	const op = "repositories.postgres.SetActive"

	conn := s.getter.DefaultTrOrDB(ctx, s.pool)

	// Получаем числовой id
	id, err := s.getUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	// Обновляем is_active у пользователя
	_, err = conn.Exec(
		ctx,
		`UPDATE users SET is_active = $1 WHERE id = $2`,
		isActive, id,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Возвращает ID (int64) пользователя по его ID (string)
func (s *Storage) getUserID(
	ctx context.Context,
	userID string,
) (int64, error) {
	conn := s.getter.DefaultTrOrDB(ctx, s.pool)
	res := conn.QueryRow(
		ctx,
		`
		SELECT u.id 
		FROM users u 
		JOIN users_id i ON u.user_id = i.id
		WHERE i.user_id = $1;
		`,
		userID,
	)

	var id int64
	err := res.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
