package store

import "github.com/jackc/pgx/v5/pgxpool"

type Store struct {
	Pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{
		Pool: pool,
	}
}
