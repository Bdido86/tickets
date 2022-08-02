package repository

import (
	"context"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/models"
)

type Cinema interface {
	GetFilms(ctx context.Context, limit uint64, offset uint64, desc bool) ([]models.Film, error)
	GetFilmRoom(ctx context.Context, filmId uint, currentUserId uint) (models.FilmRoom, error)

	AuthUser(ctx context.Context, name string) (models.User, error)
	GetUserIdByToken(ctx context.Context, token string) (uint, error)

	GetMyTickets(ctx context.Context, currentUserId uint) ([]models.Ticket, error)
	CreateTicket(ctx context.Context, filmId uint, placeId uint, currentUserId uint) (models.Ticket, error)
	DeleteTicket(ctx context.Context, ticketId uint, currentUserId uint) error
}
