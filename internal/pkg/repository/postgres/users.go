package repository

import (
	"context"
	"encoding/base64"
	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/models"
	"go.opencensus.io/trace"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type User interface {
	AuthUser(ctx context.Context, name string) (models.User, error)
	GetUserIdByToken(ctx context.Context, token string) (uint, error)
}

func (r *Repository) AuthUser(ctx context.Context, name string) (models.User, error) {
	ctx, span := trace.StartSpan(ctx, "repository/AuthUser")
	defer span.End()

	r.mu.Lock()
	defer r.mu.Unlock()

	var user models.User
	query, args, err := squirrel.Select("*").
		From("users").
		Where(
			squirrel.Eq{"name": strings.ToLower(name)},
		).PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		r.logger.Errorf("Repository.AuthUser.Select %v", err)
		return user, errors.Wrap(err, "Repository.AuthUser.Select")
	}

	if err := pgxscan.Get(ctx, r.pool, &user, query, args...); err != nil {
		if !pgxscan.NotFound(err) {
			r.logger.Errorf("Repository.AuthUser.Get %v", err)
			return user, errors.Wrap(err, "Repository.AuthUser.Get")
		}

		user, err = r.createUser(ctx, name)
		if err != nil {
			r.logger.Errorf("Repository.AuthUser error %v", err)
			return user, err
		}
	}

	return user, nil
}

func (r *Repository) GetUserIdByToken(ctx context.Context, token string) (uint, error) {
	ctx, span := trace.StartSpan(ctx, "repository/GetUserIdByToken")
	defer span.End()

	r.mu.Lock()
	defer r.mu.Unlock()

	var userId uint
	query, args, err := squirrel.Select("id").
		From("users").
		Where(
			squirrel.Eq{"token": token},
		).PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		r.logger.Errorf("Repository.GetUserIdByToken.Select %v", err)
		return userId, errors.Wrap(err, "Repository.GetUserIdByToken.Select")
	}

	if err := pgxscan.Get(ctx, r.pool, &userId, query, args...); err != nil {
		if pgxscan.NotFound(err) {
			r.logger.Errorf("Invalid Token. User not found by token %v", err)
			return userId, errors.New("Invalid Token. User not found by token")
		}
		r.logger.Errorf("Repository.GetUserIdByToken.Get %v", err)
		return userId, errors.Wrap(err, "Repository.GetUserIdByToken.Get")
	}

	return userId, nil
}

func (r *Repository) createUser(ctx context.Context, name string) (models.User, error) {
	ctx, span := trace.StartSpan(ctx, "repository/createUser")
	defer span.End()

	var user models.User
	token := generateToken(name)
	query, args, err := squirrel.Insert("users").
		Columns("name, token").
		Values(strings.ToLower(name), token).
		Suffix("RETURNING id, name, token").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		r.logger.Errorf("Repository.createUser: to sql %v", err)
		return user, errors.Wrap(err, "Repository.createUser: to sql")
	}

	row := r.pool.QueryRow(ctx, query, args...)
	if err := row.Scan(&user.Id, &user.Name, &user.Token); err != nil {
		r.logger.Errorf("Repository.createUser: insert %v", err)
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
