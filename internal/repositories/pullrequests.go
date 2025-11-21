package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/iskanye/avito-tech-internship/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Добавляет PR в базу данных
func (s Storage) CreatePullRequest(
	ctx context.Context,
	pullRequestID string,
	pullRequestName string,
	authorID string,
) error {
	const op = "repositories.postgres.CreatePullRequest"

	// Получаем ID автора
	id, err := s.getUserID(ctx, authorID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	// Вставляем ID пул реквеста
	insertID := s.pool.QueryRow(
		ctx,
		"INSERT INTO pull_requests_id (pull_request_id) VALUES ($1) RETURNING id;",
		pullRequestID,
	)

	var dbID int64
	err = insertID.Scan(&dbID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == UNIQUE_VIOLATION_CODE {
			return ErrPRExists
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	// Вставляем сам пул реквест
	insertPR := s.pool.QueryRow(
		ctx,
		`
		INSERT INTO pull_requests (
		pull_request_id, pull_request_name, author_id, status, created_at
		) VALUES ($1, $2, $3, $4, $5) RETURNING id;
		`,
		dbID, pullRequestName, id, models.PULLREQUEST_OPEN, time.Now(),
	)

	// Получаем ID пулреквеста в базе данных
	var prID int64
	err = insertPR.Scan(&prID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
