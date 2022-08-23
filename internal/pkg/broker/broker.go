//go:generate mockgen  -destination=./mocks/broker.go -package=mock_broker -source=./broker.go Broker

package broker

import (
	"context"
)

type Broker interface {
	CreateTicket(ctx context.Context, filmId, placeId uint) error
	DeleteTicket(ctx context.Context, ticketId uint) error
	Close(ctx context.Context) error
}
