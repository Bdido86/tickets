package repository

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/models"
)

func (r *Repository) GetMyTickets(ctx context.Context, currentUserId uint) ([]models.Ticket, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var tickets []models.Ticket
	query, args, err := squirrel.Select("*").
		From("tickets").
		Where(
			squirrel.Eq{"user_id": currentUserId},
		).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return tickets, errors.Wrap(err, "Repository.GetMyTickets.Select")
	}

	if err := pgxscan.Select(ctx, r.pool, &tickets, query, args...); err != nil {
		return tickets, errors.Wrap(err, "Repository.GetMyTickets.Select: error scan")
	}

	return tickets, nil
}
