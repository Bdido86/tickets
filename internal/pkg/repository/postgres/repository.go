package repository

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/logger"
	"sync"
)

type Repository struct {
	pool   *pgxpool.Pool
	mu     sync.RWMutex
	logger logger.Logger
}

func NewRepository(pool *pgxpool.Pool, logger logger.Logger) *Repository {
	return &Repository{
		pool:   pool,
		mu:     sync.RWMutex{},
		logger: logger,
	}
}
