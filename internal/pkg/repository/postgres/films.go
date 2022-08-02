package repository

import (
	"context"
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
