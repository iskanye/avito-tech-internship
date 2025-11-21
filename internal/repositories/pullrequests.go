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
		FROM pull_requests 
		WHERE id = $1; 
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

	// Получаем ревьюверов
	getReviewers, err := s.pool.Query(
		ctx,
		`
		SELECT i.user_id 
		FROM reviewers r
		JOIN users_id i ON r.user_id = i.id
		WHERE pull_request_id = $1;
		`,
		prID,
	)
	// Если ревьюверы не нашлись, то значит их и нет
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return models.PullRequest{}, fmt.Errorf("%s: %w", op, err)
	}

	for getReviewers.Next() {
		var reviewer string
		err := getReviewers.Scan(&getReviewers)
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

	// Получаем ID юзера
	id, err := s.getUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Получаем все пул реквесты в которых данный юзер ревьювер
	getPullRequests, err := s.pool.Query(
		ctx,
		`
		SELECT ip.pull_request_id, p.pull_request_name, iu.user_id, p.status
		FROM reviewers r
		JOIN pull_requests p ON p.id = r.pull_request_id
		JOIN pull_requests_id ip ON ip.id = r.pull_request_id
		JOIN users_id iu ON iu.id = p.author_id
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
) error {
	const op = "repositories.postgres.MergePullRequest"

	// Обновляем статус пул реквеста
	_, err := s.pool.Exec(
		ctx,
		`
		UPDATE pull_requests 
		SET status = $1 
		WHERE id = (
			SELECT id 
			FROM pull_requests_id 
			WHERE pull_request_id = $2
		);
		`,
		models.PULLREQUEST_MERGED, pullRequestID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
