package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	pool *pgxpool.Pool
	tx   pgx.Tx
}

func New(
	host string,
	port int,
	user string,
	password string,
	dbName string,
	maxConns int32,
) (*Storage, error) {
	const op = "repositories.postgres.New"

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable&pool_max_conns=%d",
		user, password, host, port, dbName, maxConns,
	)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		pool: pool,
	}, nil
}

func (s *Storage) Stop() {
	s.pool.Close()
}

// Начинает транзакцию
func (s *Storage) Begin(c context.Context) error {
	const op = "repositories.postgres.Begin"

	tx, err := s.pool.Begin(c)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	s.tx = tx
	return nil
}

// Отменяет изменения транзакции
func (s *Storage) Rollback(c context.Context) error {
	return s.tx.Rollback(c)
}

// Применяет изменения транзакции
func (s *Storage) Commit(c context.Context) error {
	const op = "repositories.postgres.Commit"

	err := s.tx.Commit(c)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
