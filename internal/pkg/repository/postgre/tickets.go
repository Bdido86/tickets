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

func (r *Repository) DeleteTicket(ctx context.Context, ticketId uint, currentUserId uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	query, args, err := squirrel.Select("id").
		From("tickets").
		Where(
			squirrel.Eq{"user_id": currentUserId, "id": ticketId},
		).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return errors.Wrap(err, "Repository.DeleteTicket.Select")
	}

	var id uint
	if err := pgxscan.Get(ctx, r.pool, &id, query, args...); err != nil {
		if pgxscan.NotFound(err) {
			return errors.Wrap(err, "Repository.DeleteTicket.Get: Ticket not found")
		}
		return errors.Wrap(err, "Repository.DeleteTicket.Select: error scan")
	}

	query, args, err = squirrel.Delete("tickets").
		Where(
			squirrel.Eq{"user_id": currentUserId, "id": ticketId},
		).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return errors.Wrap(err, "Repository.DeleteTicket.Delete")
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "Repository.DeleteTicket.Delete")
	}

	return nil
}
