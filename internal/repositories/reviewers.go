package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// Назначает наблюдателей на пул реквест
func (s *Storage) AssignReviewers(
	ctx context.Context,
	pullRequestID string,
	authorID string,
) error {
	const op = "repositories.postgres.AssignReviewers"

	// Получаем ID автора
	id, err := s.getUserID(ctx, authorID)
	if err != nil {
		// Ловить ErrNotFound нет смысла так как
		// при создании PR уже на это проверялось
		return fmt.Errorf("%s: %w", op, err)
	}

	// Получаем ID пул реквеста
	getID := s.pool.QueryRow(
		ctx,
		`
		SELECT p.id FROM pull_requests p
		JOIN pull_requests_id i ON i.id = p.id
		WHERE i.pull_request_id = $1;
		`,
		pullRequestID,
	)

	var prID int64
	err = getID.Scan(&prID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	// Получаем ID доступных членов команды
	getReviewers, err := s.pool.Query(
		ctx,
		`
		SELECT id FROM users
		WHERE 
			id <> $1 AND 
			team_id = (SELECT team_id from users WHERE id = $2) AND 
			is_active = TRUE
		LIMIT 2;
		`,
		id, authorID,
	)
	if err != nil {
		// Если юзеров в команде кроме самого автора нет
		// то ничего не делаем
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	defer getReviewers.Close()

	for getReviewers.Next() {
		var reviewerID int64
		err := getReviewers.Scan(&reviewerID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		// Назначаем ревьювера
		_, err = s.pool.Exec(
			ctx,
			`
			INSERT INTO reviewers (pull_request_id, user_id)
			VALUES ($1, $2);
			`,
			prID, reviewerID,
		)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

// Переназначает наблюдателя на пул реквест
func (s *Storage) ReassignReviewer(
	ctx context.Context,
	pullRequestID string,
	oldReviewerID string,
) (string, error) {
	const op = "repositories.postgres.ReassignReviewer"

	// Получаем ID и автора пул реквеста
	getID := s.pool.QueryRow(
		ctx,
		"SELECT id, author_id FROM pull_requests_id WHERE pull_request_id = $1",
		pullRequestID,
	)

	var prID, authorID int64
	err := getID.Scan(&prID, &authorID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// Получаем ID прошлого ревьювера
	getOldReviewer := s.pool.QueryRow(
		ctx,
		"SELECT id FROM user_id WHERE user_id = $1",
		oldReviewerID,
	)

	var oldReviewer int64
	err = getOldReviewer.Scan(&oldReviewer)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// Получаем нового ревьювера
	getNewReviewer := s.pool.QueryRow(
		ctx,
		`
		SELECT u.id FROM users u
		WHERE 
			u.team_id = (SELECT team_id FROM users WHERE id = $1) AND
			u.id <> $2 AND
			u.id NOT IN (
				SELECT user_id FROM reviewers
				WHERE pull_request_id = $3
			);
		`,
		oldReviewerID, authorID, pullRequestID,
	)

	var newReviewer int64
	err = getNewReviewer.Scan(&newReviewer)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// Обновляем старого ревьювера
	_, err = s.pool.Exec(
		ctx,
		`
		UPDATE reviewers SET user_id = $1
		WHERE pull_request_id = $2 AND user_id = $3
		`,
		newReviewer, prID, oldReviewer,
	)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// Получаем ID нового ревьювера
	getNewReviewerID := s.pool.QueryRow(
		ctx,
		"SELECT user_id FROM users_id WHERE id = $1",
		newReviewer,
	)

	var newReviewerID string
	err = getNewReviewerID.Scan(&getNewReviewerID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return newReviewerID, nil
}
