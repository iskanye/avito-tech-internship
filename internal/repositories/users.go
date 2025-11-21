package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/iskanye/avito-tech-internship/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Добавляет пользователя в БД
func (s *Storage) AddUser(
	ctx context.Context,
	user models.User,
) error {
	const op = "repositories.postgres.AddUser"

	// Вставляем юзера
	getUserID := s.pool.QueryRow(
		ctx,
		"INSERT INTO users (username, team_id, is_active) VALUES ($1, $2, $3) RETURNING id;",
		user.Username, user.TeamID, user.IsActive,
	)

	var id int64
	err := getUserID.Scan(&id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Вставить ID в базу
	_, err = s.pool.Exec(
		ctx,
		"INSERT INTO users_id (id, user_id) VALUES ($1, $2);",
		id, user.UserID,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == UNIQUE_VIOLATION_CODE {
			return ErrUserExists
		}
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

	// Получаем ID пользователя
	id, err := s.getUserID(ctx, user.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.pool.Exec(
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

	id, err := s.getUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, ErrNotFound
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	// Получаем данные о пользователе из БД
	res := s.pool.QueryRow(
		ctx,
		"SELECT username, team_id, is_active FROM users WHERE id = $1",
		id,
	)

	user := models.User{
		UserID: userID,
	}
	err = res.Scan(&user.Username, &user.TeamID, &user.IsActive)
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

	// Получаем числовой id
	id, err := s.getUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	// Обновляем is_active у пользователя
	_, err = s.pool.Exec(
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
	res := s.pool.QueryRow(
		ctx,
		"SELECT u.id FROM users_id u WHERE u.user_id = $1",
		userID,
	)

	var id int64
	err := res.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
