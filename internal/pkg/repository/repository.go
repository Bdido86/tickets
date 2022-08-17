//go:generate mockgen -source ./repository.go -destination=./mocks/repository.go -package=mock_repository

package repository

import (
	repository "gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/repository/postgres"
)

type Cinema interface {
	repository.Film
	repository.User
	repository.Ticket
}
