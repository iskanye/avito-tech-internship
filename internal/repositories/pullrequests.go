package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/iskanye/avito-tech-internship/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Добавляет PR в базу данных
func (s *Storage) CreatePullRequest(
	ctx context.Context,
	pullRequest models.PullRequest,
) error {
	const op = "repositories.postgres.CreatePullRequest"

	// Получаем ID автора
	id, err := s.getUserID(ctx, pullRequest.AuthorID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	// Вставляем пул реквест
	insertPR := s.pool.QueryRow(
		ctx,
		`
		INSERT INTO pull_requests (
		pull_request_name, author_id, status, created_at, merged_at
		) VALUES ($1, $2, $3, $4) RETURNING id;
		`,
		pullRequest.Name, id, pullRequest.Status, pullRequest.CreatedAt, pullRequest.MergedAt,
	)

	// Получаем ID пулреквеста в базе данных
	var prID int64
	err = insertPR.Scan(&prID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Вставляем ID пул реквеста
	_, err = s.pool.Exec(
		ctx,
		"INSERT INTO pull_requests_id (id, pull_request_id) VALUES ($1, $2)",
		prID, pullRequest.ID,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == UNIQUE_VIOLATION_CODE {
			return ErrPRExists
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Получает PR из базы данных
func (s *Storage) GetPullRequest(
	ctx context.Context,
	pullRequestID string,
) (models.PullRequest, error) {
	const op = "repositories.postgres.GetPullRequest"

	// Получаем ID пул реквеста
	getID := s.pool.QueryRow(
		ctx,
		"SELECT id FROM pull_requests_id WHERE pull_request_id = $1;",
		pullRequestID,
	)

	var prID int64
	err := getID.Scan(&prID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.PullRequest{}, ErrNotFound
		}
		return models.PullRequest{}, fmt.Errorf("%s: %w", op, err)
	}

	// Получаем пул реквест
	getPR := s.pool.QueryRow(
		ctx,
		`
		SELECT pull_request_name, author_id, status, created_at, merged_at
		FROM pull_requests WHERE id = $1; 
		`,
		prID,
	)

	var pullRequest models.PullRequest
	var author string
	err = getPR.Scan(
		&pullRequest.Name,
		&author,
		&pullRequest.Status,
		&pullRequest.CreatedAt,
		&pullRequest.MergedAt,
	)
	if err != nil {
		// Нет смысла проверять на ErrNoRows
		return models.PullRequest{}, fmt.Errorf("%s: %w", op, err)
	}

	// Получаем ID автора
	getAuthorID := s.pool.QueryRow(
		ctx,
		`
		SELECT i.user_id 
		FROM users u
		JOIN users_id i ON i.id = u.user_id
		WHERE u.id = $1;
		`,
		author,
	)

	err = getAuthorID.Scan(&pullRequest.AuthorID)
	if err != nil {
		// Нет смысла проверять на ErrNoRows
		return models.PullRequest{}, fmt.Errorf("%s: %w", op, err)
	}

	return pullRequest, nil
}
