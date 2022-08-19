//go:build integration
// +build integration

package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	internalModels "gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/models"
	"gitlab.ozon.dev/Bdido86/movie-tickets/tests/postgres"
	"testing"
)

func TestGetMyTickets(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Run("empty", func(t *testing.T) {
			Db.SetUp(t)
			defer Db.TearDown()

			//arrange
			userId := uint(1)
			var apiResp []internalModels.Ticket

			// act
			res, err := Repository.GetMyTickets(context.Background(), userId)

			// assert
			assert.NoError(t, err)
			assert.Equal(t, res, apiResp)
		})

		t.Run("list", func(t *testing.T) {
			//arrange
			Db.SetUp(t)
			defer Db.TearDown()

			ctx := context.Background()
			userId := uint(1)

			// insert tickets
			query := postgres.InsertTicket
			var args []interface{}
			args = append(args, userId) // user_id
			args = append(args, 1)      // film_id
			args = append(args, 1)      // room_id
			args = append(args, 5)      // place
			Db.Pool.QueryRow(ctx, query, args...)

			args = nil
			args = append(args, userId) // user_id
			args = append(args, 2)      // film_id
			args = append(args, 2)      // room_id
			args = append(args, 3)      // place
			Db.Pool.QueryRow(ctx, query, args...)

			apiResp := []internalModels.Ticket{
				{
					Id:     uint64(1),
					UserId: uint64(1),
					FilmId: uint64(1),
					RoomId: uint64(1),
					Place:  uint64(5),
				},
				{
					Id:     uint64(2),
					UserId: uint64(userId),
					FilmId: uint64(2),
					RoomId: uint64(2),
					Place:  uint64(3),
				},
			}

			// act
			res, err := Repository.GetMyTickets(context.Background(), userId)

			// assert
			assert.NoError(t, err)
			assert.Equal(t, res, apiResp)
		})
	})
}
