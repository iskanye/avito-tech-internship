package repositories

import (
	"context"
	"fmt"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func New(
	host string,
	port int,
	user string,
	password string,
	dbName string,
	maxConns int32,
	getter *trmpgx.CtxGetter,
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
		pool:   pool,
		getter: getter,
	}, nil
}

func (s *Storage) Stop() {
	s.pool.Close()
}

func (s *Storage) GetPool() *pgxpool.Pool {
	return s.pool
}
