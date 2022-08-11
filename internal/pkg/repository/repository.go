package repository

import (
	repository "gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/repository/postgre"
)

type Cinema interface {
	repository.Film
	repository.User
	repository.Ticket
}
