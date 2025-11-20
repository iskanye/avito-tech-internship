package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	pool *pgxpool.Pool
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
		"postges://%s:%s@%s:%d/%s",
		user, password, host, port, dbName,
	)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	config.MaxConns = maxConns

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
