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
func (s *Storage) CreatePullRequest(
	ctx context.Context,
	pullRequest models.PullRequest,
) error {
	const op = "repositories.postgres.CreatePullRequest"

	conn := s.getter.DefaultTrOrDB(ctx, s.pool)

	// Получаем ID автора
	id, err := s.getUserID(ctx, pullRequest.AuthorID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	// Вставляем ID пул реквеста
	insertID := conn.QueryRow(
		ctx,
		"INSERT INTO pull_requests_id (pull_request_id) VALUES ($1) RETURNING id",
		pullRequest.ID,
	)

	var prID int64
	err = insertID.Scan(&prID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == UNIQUE_VIOLATION_CODE {
			return ErrPRExists
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	// Вставляем пул реквест
	_, err = conn.Exec(
		ctx,
		`
		INSERT INTO pull_requests (
			pull_request_id, pull_request_name, author_id, status, created_at, merged_at
		) 
		VALUES ($1, $2, $3, $4, $5, $6);
		`,
		prID, pullRequest.Name, id, pullRequest.Status, pullRequest.CreatedAt, pullRequest.MergedAt,
	)
	if err != nil {
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

	conn := s.getter.DefaultTrOrDB(ctx, s.pool)

	// Получаем пул реквест
	getPR := conn.QueryRow(
		ctx,
		`
		SELECT p.id, p.pull_request_name, p.author_id, p.status, p.created_at, p.merged_at 
		FROM pull_requests p 
		JOIN pull_requests_id i ON p.pull_request_id = i.id
		WHERE i.pull_request_id = $1;
		`,
		pullRequestID,
	)

	pullRequest := models.PullRequest{
		ID: pullRequestID,
	}
	var author, prID int64

	err := getPR.Scan(
		&prID,
		&pullRequest.Name,
		&author,
		&pullRequest.Status,
		&pullRequest.CreatedAt,
		&pullRequest.MergedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.PullRequest{}, ErrNotFound
		}
		return models.PullRequest{}, fmt.Errorf("%s: %w", op, err)
	}

	// Получаем ID автора
	getAuthorID := conn.QueryRow(
		ctx,
		`
		SELECT i.user_id 
		FROM users u
		JOIN users_id i ON u.user_id = i.id
		WHERE u.id = $1;
		`,
		author,
	)

	err = getAuthorID.Scan(&pullRequest.AuthorID)
	if err != nil {
		// Нет смысла проверять на ErrNoRows
		return models.PullRequest{}, fmt.Errorf("%s: %w", op, err)
	}

	// Получаем ревьюверов
	getReviewers, err := conn.Query(
		ctx,
		`
		SELECT i.user_id 
		FROM reviewers r
		JOIN users u ON r.user_id = u.id
		JOIN users_id i ON u.user_id = i.id
		WHERE pull_request_id = $1;
		`,
		prID,
	)
	if err != nil {
		// Если ревьюверы не нашлись, то значит их и нет
		if !errors.Is(err, pgx.ErrNoRows) {
			pullRequest.AssignedReviewers = []string{}
			return pullRequest, nil
		}
		return models.PullRequest{}, fmt.Errorf("%s: %w", op, err)
	}

	for getReviewers.Next() {
		var reviewer string
		err := getReviewers.Scan(&reviewer)
		if err != nil {
			return models.PullRequest{}, fmt.Errorf("%s: %w", op, err)
		}

		pullRequest.AssignedReviewers = append(pullRequest.AssignedReviewers, reviewer)
	}

	return pullRequest, nil
}

// Получает PR в который данные пользователь назначен ревьювером
func (s *Storage) GetReview(
	ctx context.Context,
	userID string,
) ([]models.PullRequest, error) {
	const op = "repositories.postgres.GetReview"

	conn := s.getter.DefaultTrOrDB(ctx, s.pool)

	// Получаем ID юзера
	id, err := s.getUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Получаем все пул реквесты в которых данный юзер ревьювер
	getPullRequests, err := conn.Query(
		ctx,
		`
		SELECT ip.pull_request_id, p.pull_request_name, iu.user_id, p.status
		FROM reviewers r
		JOIN pull_requests p ON r.pull_request_id = p.id
		JOIN pull_requests_id ip ON ip.id = p.pull_request_id
		JOIN users u ON p.author_id = u.id
		JOIN users_id iu ON u.user_id = iu.id
		WHERE r.user_id = $1;
		`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer getPullRequests.Close()

	// Читаем строчки
	var pullRequests []models.PullRequest
	for getPullRequests.Next() {
		var pullRequest models.PullRequest
		err := getPullRequests.Scan(
			&pullRequest.ID,
			&pullRequest.Name,
			&pullRequest.AuthorID,
			&pullRequest.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		pullRequests = append(pullRequests, pullRequest)
	}

	return pullRequests, nil
}

// Помечает PR как MERGED
func (s *Storage) MergePullRequest(
	ctx context.Context,
	pullRequestID string,
	mergedAt time.Time,
) error {
	const op = "repositories.postgres.MergePullRequest"

	conn := s.getter.DefaultTrOrDB(ctx, s.pool)

	// Обновляем статус пул реквеста
	_, err := conn.Exec(
		ctx,
		`
		UPDATE pull_requests 
		SET status = $1, merged_at = $2 
		WHERE pull_request_id = (
			SELECT id 
			FROM pull_requests_id 
			WHERE pull_request_id = $3
		);
		`,
		models.PULLREQUEST_MERGED, mergedAt, pullRequestID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
