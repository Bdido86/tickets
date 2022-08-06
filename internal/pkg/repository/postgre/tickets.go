package repository

import (
	"context"
	"fmt"
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

func (r *Repository) CreateTicket(ctx context.Context, filmId uint, placeId uint, currentUserId uint) (models.Ticket, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var ticket models.Ticket

	query, args, err := squirrel.Select("*").
		From("films").
		Where(
			squirrel.Eq{"id": filmId},
		).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return ticket, errors.Wrap(err, "Repository.CreateTicket.SelectFilm")
	}
	var film models.Film
	if err := pgxscan.Get(ctx, r.pool, &film, query, args...); err != nil {
		return ticket, errors.Wrap(err, "Repository.CreateTicket.SelectFilm: film not found")
	}

	query, args, err = squirrel.Select("rooms.*").
		From("film_room").
		Join("rooms ON film_room.room_id=rooms.id").
		Where(
			squirrel.Eq{"film_room.film_id": film.Id},
		).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return ticket, errors.Wrap(err, "Repository.CreateTicket.SelectFilmRoom query")
	}
	var roomDb models.RoomDb
	if err := pgxscan.Get(ctx, r.pool, &roomDb, query, args...); err != nil {
		if pgxscan.NotFound(err) {
			return ticket, errors.Wrap(err, "Repository.CreateTicket.SelectFilmRoom: film_room not found")
		}
		return ticket, errors.Wrap(err, "Repository.CreateTicket.SelectFilmRoom")
	}
	if placeId > uint(roomDb.CountPlaces) {
		return ticket, errors.New("Repository.CreateTicket.SelectFilmRoom: place not exist in room")
	}

	query, args, err = squirrel.Select("*").
		From("tickets").
		Where(
			squirrel.Eq{"film_id": film.Id, "room_id": roomDb.Id, "place": placeId},
		).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return ticket, errors.Wrap(err, "Repository.CreateTicket.SelectTickets")
	}

	if err := pgxscan.Get(ctx, r.pool, &ticket, query, args...); err != nil {
		if !pgxscan.NotFound(err) {
			return ticket, errors.Wrap(err, "Repository.CreateTicket.SelectTickets")
		}
		fmt.Printf("%+v\n", ticket)
		query, args, err = squirrel.Insert("tickets").
			Columns("user_id, film_id, room_id, place").
			Values(currentUserId, filmId, roomDb.Id, placeId).
			Suffix("RETURNING id, user_id, film_id, room_id, place").
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if err != nil {
			return ticket, errors.Wrap(err, "Repository.createUser: to sql")
		}
		row := r.pool.QueryRow(ctx, query, args...)
		if err := row.Scan(&ticket.Id, &ticket.UserId, &ticket.FilmId, &ticket.RoomId, &ticket.Place); err != nil {
			return ticket, errors.Wrap(err, "Repository.CreateTicket.Insert")
		}
		return ticket, nil
	}

	return models.Ticket{}, errors.New("Repository.CreateTicket: seat is taken")
}
