package repository

import (
	"context"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/models"
)

type Cinema interface {
	GetFilms(ctx context.Context) ([]models.Film, error)
	AuthUser(ctx context.Context, name string) (models.User, error)
}
