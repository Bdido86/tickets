package repository

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"sync"
)

type Repository struct {
	pool *pgxpool.Pool
	mu   sync.RWMutex
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
		mu:   sync.RWMutex{},
	}
}
