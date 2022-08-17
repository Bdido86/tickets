package repository

import (
	"context"
	"encoding/base64"
	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type User interface {
	AuthUser(ctx context.Context, name string) (models.User, error)
	GetUserIdByToken(ctx context.Context, token string) (uint, error)
}

func (r *Repository) AuthUser(ctx context.Context, name string) (models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var user models.User
	query, args, err := squirrel.Select("*").
		From("users").
		Where(
			squirrel.Eq{"name": strings.ToLower(name)},
		).PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		return user, errors.Wrap(err, "Repository.AuthUser.Select")
	}

	if err := pgxscan.Get(ctx, r.pool, &user, query, args...); err != nil {
		if !pgxscan.NotFound(err) {
			return user, errors.Wrap(err, "Repository.AuthUser.Get")
		}

		user, err = r.createUser(ctx, name)
		if err != nil {
			return user, err
		}
	}

	return user, nil
}

func (r *Repository) GetUserIdByToken(ctx context.Context, token string) (uint, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var userId uint
	query, args, err := squirrel.Select("id").
		From("users").
		Where(
			squirrel.Eq{"token": token},
		).PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		return userId, errors.Wrap(err, "Repository.GetUserIdByToken.Select")
	}

	if err := pgxscan.Get(ctx, r.pool, &userId, query, args...); err != nil {
		if pgxscan.NotFound(err) {
			return userId, errors.New("Invalid Token. User not found by token")
		}
		return userId, errors.Wrap(err, "Repository.GetUserIdByToken.Get")
	}

	return userId, nil
}

func (r *Repository) createUser(ctx context.Context, name string) (models.User, error) {
	var user models.User
	token := generateToken(name)
	query, args, err := squirrel.Insert("users").
		Columns("name, token").
		Values(strings.ToLower(name), token).
		Suffix("RETURNING id, name, token").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return user, errors.Wrap(err, "Repository.createUser: to sql")
	}

	row := r.pool.QueryRow(ctx, query, args...)
	if err := row.Scan(&user.Id, &user.Name, &user.Token); err != nil {
		return user, errors.Wrap(err, "Repository.createUser: insert")
	}
	return user, nil
}

func generateToken(name string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(name), bcrypt.DefaultCost)
	if err != nil {
		return name
	}

	return base64.StdEncoding.EncodeToString(hash)
}
