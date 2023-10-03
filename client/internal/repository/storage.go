package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Storage struct {
	pgConnPool *pgxpool.Pool
}

func NewStorage(connPool *pgxpool.Pool) *Storage {
	return &Storage{
		pgConnPool: connPool,
	}
}

func (s *Storage) GetPgConnPool() *pgxpool.Pool {
	return s.pgConnPool
}
