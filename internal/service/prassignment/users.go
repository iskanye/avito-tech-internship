package prassignment

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/iskanye/avito-tech-internship/internal/models"
	"github.com/iskanye/avito-tech-internship/internal/repositories"
)

// Назначает пользователю свойство is_active
func (a *PRAssignment) SetIsActive(
	ctx context.Context,
	userID string,
	isActive bool,
) (models.User, error) {
	const op = "service.PRAssignment.SetIsActive"

	log := a.log.With(
		slog.String("op", op),
		slog.String("user_id", userID),
	)

	log.Info("Attempting to set is_active")

	// Начинаем транзакцию
	err := a.txManager.Begin(ctx)
	if err != nil {
		log.Error("Failed to start transaction",
			slog.String("err", err.Error()),
		)
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	defer a.txManager.Rollback(ctx)

	// Обновляем is_active пользователя
	err = a.userModifier.SetActive(ctx, userID, isActive)
	if err != nil {
		log.Error("Failed to set is_active",
			slog.String("err", err.Error()),
		)
		if errors.Is(err, repositories.ErrNotFound) {
			return models.User{}, ErrNotFound
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	// Сохраняем изменения
	if err = a.txManager.Commit(ctx); err != nil {
		log.Error("Failed to commit transaction",
			slog.String("err", err.Error()),
		)
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	// Получаем пользователя, чтобы вернуть
	user, err := a.userProvider.GetUser(ctx, userID)
	// Если на прошлом этапе уже не вылетела ошибка ErrNotFound
	// то тут уже нет смысла её вылавливать
	if err != nil {
		log.Error("Failed to get user",
			slog.String("err", err.Error()),
		)
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Set is_active successfully")

	return user, nil
}

// Получает пул реквесты, в которых пользователь - ревьювер
func (a *PRAssignment) GetReview(
	ctx context.Context,
	userID string,
) ([]models.PullRequest, error) {
	const op = "service.PRAssignment.GetReview"

	log := a.log.With(
		slog.String("op", op),
		slog.String("user_id", userID),
	)

	log.Info("Attempting to get reviewing PRs")

	// Получаем пул реквесты
	pullRequests, err := a.prProvider.GetReview(ctx, userID)
	if err != nil {
		log.Error("Failed to get PRs",
			slog.String("err", err.Error()),
		)
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Got reviewing PRs successfully")

	return pullRequests, nil
}
