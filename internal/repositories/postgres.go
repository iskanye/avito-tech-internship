package repositories

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	db *sql.DB
}

func New(
	host string,
	port int,
	user string,
	password string,
	dbName string,
) (*Storage, error) {
	const op = "repositories.postgres.New"

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName,
	)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) Stop() error {
	return s.db.Close()
}
