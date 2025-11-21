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

	// Обновляем is_active пользователя
	err := a.userModifier.SetActive(ctx, userID, isActive)
	if err != nil {
		log.Error("Failed to set is_active",
			slog.String("err", err.Error()),
		)
		if errors.Is(err, repositories.ErrNotFound) {
			return models.User{}, ErrNotFound
		}

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
