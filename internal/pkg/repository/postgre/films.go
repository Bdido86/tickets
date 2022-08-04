package repository

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/models"
)

func (r *Repository) GetFilms(ctx context.Context, limit uint64, offset uint64, desc bool) ([]models.Film, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	builder := squirrel.Select("*").From("films")
	if limit > 0 {
		builder = builder.Limit(limit)
	}
	if offset > 0 {
		builder = builder.Offset(offset)
	}
	if desc {
		builder = builder.OrderBy("name DESC")
	} else {
		builder = builder.OrderBy("name ASC")
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Repository.GetFilms.ToSql")
	}

	var films []models.Film
	if err := pgxscan.Select(ctx, r.pool, &films, query, args...); err != nil {
		if pgxscan.NotFound(err) {
			return films, nil
		}
		return nil, errors.Wrap(err, "Repository.GetFilms.Select")
	}

	return films, nil
}

func (r *Repository) GetFilmRoom(ctx context.Context, filmId uint, currentUserId uint) (models.FilmRoom, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var filmRoom models.FilmRoom

	query, args, err := squirrel.Select("*").
		From("films").
		Where(
			squirrel.Eq{"id": filmId},
		).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return filmRoom, errors.Wrap(err, "Repository.GetFilms.ToSql")
	}
	var film models.Film
	if err := pgxscan.Get(ctx, r.pool, &film, query, args...); err != nil {
		if pgxscan.NotFound(err) {
			return filmRoom, errors.Wrap(err, "Repository.GetFilmRoom.Get: film not found")
		}
		return filmRoom, errors.Wrap(err, "Repository.GetFilms.Get")
	}

	query, args, err = squirrel.Select("rooms.*").
		From("film_room").
		Join("rooms ON film_room.room_id=rooms.id").
		Where(
			squirrel.Eq{"film_room.film_id": film.Id},
		).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return filmRoom, errors.Wrap(err, "Repository.GetFilms.film_room")
	}
	var roomDb models.RoomDb
	if err := pgxscan.Get(ctx, r.pool, &roomDb, query, args...); err != nil {
		if pgxscan.NotFound(err) {
			return filmRoom, errors.Wrap(err, "Repository.GetFilmRoom.Get: film_room not found")
		}
		return filmRoom, errors.Wrap(err, "Repository.GetFilms.Get")
	}

	query, args, err = squirrel.Select("*").
		From("tickets").
		Where(
			squirrel.Eq{"film_id": film.Id, "room_id": roomDb.Id},
		).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return filmRoom, errors.Wrap(err, "Repository.GetFilms.tickets")
	}
	var tickets []models.Ticket
	if err := pgxscan.Select(ctx, r.pool, &tickets, query, args...); err != nil {
		return filmRoom, errors.Wrap(err, "Repository.GetFilms.tickets: error scan")
	}

	ticketsByPlaceId := make(map[uint64]models.Ticket)
	for _, ticket := range tickets {
		ticketsByPlaceId[ticket.Place] = ticket
	}

	places := make([]models.Place, 0, roomDb.CountPlaces)
	var i uint64
	for i = 1; i <= roomDb.CountPlaces; i++ {
		var IsMy, IsFree bool
		findTicket, ok := ticketsByPlaceId[i]
		if !ok {
			IsMy = false
			IsFree = true
		} else {
			IsFree = false
			if findTicket.UserId == uint64(currentUserId) {
				IsMy = true
			}
		}

		place := models.Place{
			Id:     i,
			IsMy:   IsMy,
			IsFree: IsFree,
		}
		places = append(places, place)
	}

	room := models.Room{
		Id:     roomDb.Id,
		Places: places,
	}

	filmRoom = models.FilmRoom{
		Film: film,
		Room: room,
	}
	return filmRoom, nil
}
