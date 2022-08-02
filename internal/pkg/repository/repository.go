package repository

import (
	"context"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/models"
)

type Cinema interface {
	GetFilms(ctx context.Context) ([]models.Film, error)
	GetFilmRoom(ctx context.Context, filmId uint, currentUserId uint) (models.FilmRoom, error)

	AuthUser(ctx context.Context, name string) (models.User, error)
	GetUserIdByToken(ctx context.Context, token string) (uint, error)

	GetMyTickets(ctx context.Context, currentUserId uint) ([]models.Ticket, error)
}
