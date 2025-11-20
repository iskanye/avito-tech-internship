package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/iskanye/avito-tech-internship/internal/models"
)

// Добавляет пользователя в БД
func (s *Storage) AddUser(
	ctx context.Context,
	user models.User,
) error {
	const op = "repositories.postgres.AddUser"

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()

	// Вставить ID в базу
	insertID, err := s.db.Prepare("INSERT INTO users_id (user_id) VALUES ($1) RETURNING id;")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	res := insertID.QueryRowContext(ctx, user.UserID)

	var id int64
	err = res.Scan(&id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Вставить самого юзера
	insertUser, err := s.db.Prepare("INSERT INTO users (user_id, team_id, is_active) VALUES ($1, $2, $3);")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = insertUser.ExecContext(ctx, id, user.TeamID, user.IsActive)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = tx.Commit(); err != nil {
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

	id, err := s.getUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()

	stmt, err := s.db.Prepare(
		`UPDATE users SET is_active = $1 WHERE user_id = $2`,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, isActive, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = tx.Commit(); err != nil {
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

	stmt, err := s.db.Prepare("SELECT team_id, is_active FROM users WHERE user_id = $1")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	res := stmt.QueryRowContext(ctx, id)

	user := models.User{
		UserID: userID,
	}
	err = res.Scan(&user.TeamID, &user.IsActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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
	stmt, err := s.db.Prepare("SELECT id FROM users_id u WHERE u.user_id = $1")
	if err != nil {
		return 0, err
	}

	res := stmt.QueryRowContext(ctx, userID)

	var id int64
	err = res.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNotFound
		}
		return 0, err
	}

	return id, nil
}
