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

	conn := s.getter.DefaultTrOrDB(ctx, s.pool)

	// Получаем ID автора
	id, err := s.getUserID(ctx, authorID)
	if err != nil {
		// Ловить ErrNotFound нет смысла так как
		// при создании PR уже на это проверялось
		return fmt.Errorf("%s: %w", op, err)
	}

	// Получаем ID пул реквеста
	getID := conn.QueryRow(
		ctx,
		`
		SELECT p.id 
		FROM pull_requests p
		JOIN pull_requests_id i ON p.pull_request_id = i.id
		WHERE i.pull_request_id = $1;
		`,
		pullRequestID,
	)

	var prID int64
	err = getID.Scan(&prID)
	if err != nil {
		// Ловить ErrNotFound нет смысла так как
		// PR должен быть создан
		return fmt.Errorf("%s: %w", op, err)
	}

	// Получаем ID доступных членов команды
	getReviewers, err := conn.Query(
		ctx,
		`
		SELECT u.id 
		FROM users u
		WHERE 
			u.id <> $1 AND 
			u.team_id = (SELECT uu.team_id from users uu WHERE uu.id = $1) AND 
			u.is_active = TRUE
		ORDER BY RANDOM()
		LIMIT 2;
		`,
		id,
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

	reviewers := make([]int64, 0)
	for getReviewers.Next() {
		var reviewerID int64
		err := getReviewers.Scan(&reviewerID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		reviewers = append(reviewers, reviewerID)
	}

	for _, reviewerID := range reviewers {
		// Назначаем ревьюверов
		_, err = conn.Exec(
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

	conn := s.getter.DefaultTrOrDB(ctx, s.pool)

	// Получаем ID и автора пул реквеста
	getID := conn.QueryRow(
		ctx,
		`
		SELECT p.id, p.author_id
		FROM pull_requests p
		JOIN pull_requests_id i ON p.pull_request_id = i.id
		WHERE i.pull_request_id = $1
		`,
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
	oldReviewer, err := s.getUserID(ctx, oldReviewerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// Получаем нового ревьювера
	getNewReviewer := conn.QueryRow(
		ctx,
		`
		SELECT u.id 
		FROM users u
		WHERE 
			u.is_active = TRUE AND
			u.team_id = (
				SELECT team_id 
				FROM users 
				WHERE id = $1
			) AND
			u.id <> $2 AND
			u.id NOT IN (
				SELECT user_id 
				FROM reviewers
				WHERE pull_request_id = $3
			);
		`,
		oldReviewer, authorID, prID,
	)

	var newReviewer int64
	err = getNewReviewer.Scan(&newReviewer)
	if err != nil {
		// Нету подходящего ревьювера
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrNoCandidates
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// Обновляем старого ревьювера
	_, err = conn.Exec(
		ctx,
		`
		UPDATE reviewers 
		SET user_id = $1
		WHERE pull_request_id = $2 AND user_id = $3
		`,
		newReviewer, prID, oldReviewer,
	)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// Получаем ID нового ревьювера
	getNewReviewerID := conn.QueryRow(
		ctx,
		`
		SELECT i.user_id 
		FROM users u
		JOIN users_id i ON u.user_id = i.id 
		WHERE u.id = $1`,
		newReviewer,
	)

	var newReviewerID string
	err = getNewReviewerID.Scan(&newReviewerID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return newReviewerID, nil
}
