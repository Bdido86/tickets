package cache

import "context"
import "gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/models"

type Cache interface {
	GetUserTickets(ctx context.Context, userId uint) ([]models.Ticket, error)
	SetUserTickets(ctx context.Context, userId uint, tickets []models.Ticket) bool
	ResetUserTickets(ctx context.Context, userId uint) bool
	CounterInfo()
}
