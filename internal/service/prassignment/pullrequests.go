package prassignment

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/iskanye/avito-tech-internship/internal/models"
	"github.com/iskanye/avito-tech-internship/internal/repositories"
)

// Создаёт пул реквест
func (a *PRAssignment) CreatePullRequest(
	ctx context.Context,
	pullRequest models.PullRequest,
) (models.PullRequest, error) {
	const op = "service.PRAssignment.CreatePullRequest"

	log := a.log.With(
		slog.String("op", op),
		slog.String("pull_request_id", pullRequest.ID),
		slog.String("pull_request_name", pullRequest.Name),
	)

	log.Info("Attempting to create PR")

	// Начинаем транзакцию
	err := a.txManager.Begin(ctx)
	if err != nil {
		log.Error("Failed to start transaction",
			slog.String("err", err.Error()),
		)
		return models.PullRequest{}, fmt.Errorf("%s: %w", op, err)
	}
	defer a.txManager.Rollback(ctx)

	// Создаем пул реквест
	err = a.prCreator.CreatePullRequest(ctx, pullRequest)
	if err != nil {
		log.Error("Failed to create PR",
			slog.String("err", err.Error()),
		)
		if errors.Is(err, repositories.ErrPRExists) {
			return models.PullRequest{}, ErrPRExists
		} else if errors.Is(err, repositories.ErrNotFound) {
			return models.PullRequest{}, ErrNotFound
		}

		return models.PullRequest{}, fmt.Errorf("%s: %w", op, err)
	}

	// Назначаем ревьюверов
	err = a.revAssigner.AssignReviewers(ctx, pullRequest.ID, pullRequest.AuthorID)
	if err != nil {
		log.Error("Failed to assign reviewer",
			slog.String("err", err.Error()),
		)

		return models.PullRequest{}, fmt.Errorf("%s: %w", op, err)
	}

	// Сохраняем изменения
	if err = a.txManager.Commit(ctx); err != nil {
		log.Error("Failed to commit transaction",
			slog.String("err", err.Error()),
		)
		return models.PullRequest{}, fmt.Errorf("%s: %w", op, err)
	}

	// Получаем пул реквест
	pullRequest, err = a.prProvider.GetPullRequest(ctx, pullRequest.ID)
	if err != nil {
		// Проверять на ErrNotFound нет смысла
		log.Error("Failed to get PR",
			slog.String("err", err.Error()),
		)
		return models.PullRequest{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("PR successfully created")

	return pullRequest, nil
}

// Помечает пул реквест как MERGED
func (a *PRAssignment) MergePullRequest(
	ctx context.Context,
	pullRequestID string,
) (models.PullRequest, error) {
	const op = "service.PRAssignment.MergePullRequest"

	log := a.log.With(
		slog.String("op", op),
		slog.String("pull_request_id", pullRequestID),
	)

	log.Info("Attempting to merge PR")

	// Начинаем транзакцию
	err := a.txManager.Begin(ctx)
	if err != nil {
		log.Error("Failed to start transaction",
			slog.String("err", err.Error()),
		)
		return models.PullRequest{}, fmt.Errorf("%s: %w", op, err)
	}
	defer a.txManager.Rollback(ctx)

	// Мерджим пул реквест
	err = a.prModifier.MergePullRequest(ctx, pullRequestID)
	if err != nil {
		log.Error("Failed to merge PR",
			slog.String("err", err.Error()),
		)
		if errors.Is(err, repositories.ErrNotFound) {
			return models.PullRequest{}, ErrNotFound
		}

		return models.PullRequest{}, fmt.Errorf("%s: %w", op, err)
	}

	// Сохраняем изменения
	if err = a.txManager.Commit(ctx); err != nil {
		log.Error("Failed to commit transaction",
			slog.String("err", err.Error()),
		)
		return models.PullRequest{}, fmt.Errorf("%s: %w", op, err)
	}

	// Получаем пул реквест
	pullRequest, err := a.prProvider.GetPullRequest(ctx, pullRequestID)
	if err != nil {
		// Проверять на ErrNotFound нет смысла
		log.Error("Failed to get PR",
			slog.String("err", err.Error()),
		)
		return models.PullRequest{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("PR successfully merged")

	return pullRequest, nil
}

func (a *PRAssignment) ReassignPullRequest(
	ctx context.Context,
	pullRequestID string,
	oldReviewerID string,
) (models.PullRequest, string, error) {
	const op = "service.PRAssignment.ReassignPullRequest"

	log := a.log.With(
		slog.String("op", op),
		slog.String("pull_request_id", pullRequestID),
		slog.String("old_reviewer_id", oldReviewerID),
	)

	log.Info("Attempting to reassign PR reviewer")

	// Получаем пул реквест
	pullRequest, err := a.prProvider.GetPullRequest(ctx, pullRequestID)
	if err != nil {
		log.Error("Failed to get PR",
			slog.String("err", err.Error()),
		)
		if errors.Is(err, repositories.ErrNotFound) {
			return models.PullRequest{}, "", ErrNotFound
		}

		return models.PullRequest{}, "", fmt.Errorf("%s: %w", op, err)
	}

	// Проверяем что пул реквест не MERGED
	if pullRequest.Status == models.PULLREQUEST_MERGED {
		log.Error("Cannot reassign reviewers of merged PR")

		return models.PullRequest{}, "", ErrPRMerged
	}

	// Начинаем транзакцию
	err = a.txManager.Begin(ctx)
	if err != nil {
		log.Error("Failed to start transaction",
			slog.String("err", err.Error()),
		)
		return models.PullRequest{}, "", fmt.Errorf("%s: %w", op, err)
	}
	defer a.txManager.Rollback(ctx)

	// Переназначаем ревьювера
	newReviewerID, err := a.revModifier.ReassignReviewer(ctx, pullRequestID, oldReviewerID)
	if err != nil {
		log.Error("Failed to reassign PR reviewer",
			slog.String("err", err.Error()),
		)
		if errors.Is(err, repositories.ErrNotFound) {
			return models.PullRequest{}, "", ErrNotFound
		}

		return models.PullRequest{}, "", fmt.Errorf("%s: %w", op, err)
	}

	// Сохраняем изменения
	if err = a.txManager.Commit(ctx); err != nil {
		log.Error("Failed to commit transaction",
			slog.String("err", err.Error()),
		)
		return models.PullRequest{}, "", fmt.Errorf("%s: %w", op, err)
	}

	// Переназначаем ревьювера в нашем пул реквесте (чтобы не брать его снова из БД)
	for i, reviewers := range pullRequest.AssignedReviewers {
		if reviewers == oldReviewerID {
			pullRequest.AssignedReviewers[i] = newReviewerID
			break
		}
	}

	log.Info("Reviewer successfully reassigned")

	return pullRequest, newReviewerID, nil
}
