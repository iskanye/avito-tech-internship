package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/iskanye/avito-tech-internship/internal/models"
	"github.com/jackc/pgx/v5"
)

// Добавляет пользователя в БД
func (s *Storage) AddUser(
	ctx context.Context,
	user models.User,
) error {
	const op = "repositories.postgres.AddUser"

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback(ctx)

	// Вставить ID в базу
	insertID := s.pool.QueryRow(
		ctx,
		"INSERT INTO users_id (user_id) VALUES ($1) RETURNING id;",
		user.UserID)

	var id int64
	err = insertID.Scan(&id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Вставить самого юзера
	_, err = s.pool.Exec(
		ctx,
		"INSERT INTO users (user_id, team_id, is_active) VALUES ($1, $2, $3);",
		id, user.TeamID, user.IsActive,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
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
		return fmt.Errorf("%s: %w", op, err)
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback(ctx)

	// Обновляем is_active у пользователя
	_, err = s.pool.Exec(
		ctx,
		`UPDATE users SET is_active = $1 WHERE user_id = $2`,
		isActive, id,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Возвращает пользователя
func (s *Storage) GetUser(
	ctx context.Context,
	userID string,
) (models.User, error) {
	const op = "repositories.postgres.GetUser"

	id, err := s.getUserID(ctx, userID)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	res := s.pool.QueryRow(
		ctx,
		"SELECT team_id, is_active FROM users WHERE user_id = $1",
		id,
	)

	user := models.User{
		UserID: userID,
	}
	err = res.Scan(&user.TeamID, &user.IsActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, ErrNotFound
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// Возвращает ID (int64) пользователя по его ID (string)
func (s *Storage) getUserID(
	ctx context.Context,
	userID string,
) (int64, error) {
	res := s.pool.QueryRow(
		ctx,
		"SELECT id FROM users_id u WHERE u.user_id = $1",
		userID,
	)

	var id int64
	err := res.Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrNotFound
		}
		return 0, err
	}

	return id, nil
}
