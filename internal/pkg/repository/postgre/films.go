package repository

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/models"
)

func (r *Repository) GetFilms(ctx context.Context) ([]models.Film, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	query, args, err := squirrel.Select("*").From("films").ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "Repository.GetFilms.ToSql")
	}

	var films []models.Film
	if err := pgxscan.Select(ctx, r.pool, &films, query, args...); err != nil {
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

	fmt.Printf("%+v\n", roomDb)
	return filmRoom, nil
}
