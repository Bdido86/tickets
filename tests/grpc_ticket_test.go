//go:build integration
// +build integration

package tests

import (
	"context"
	internalModels "gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/models"
	"gitlab.ozon.dev/Bdido86/movie-tickets/pkg/api/server"
	"gitlab.ozon.dev/Bdido86/movie-tickets/tests/postgres"
	"google.golang.org/grpc/metadata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTicketCreate(t *testing.T) {
	t.Run("fail", func(t *testing.T) {
		t.Run("empty header token", func(t *testing.T) {
			//arrange
			Db.SetUp(t)
			defer Db.TearDown()

			//act
			resp, err := ServerClient.TicketCreate(context.Background(), &server.TicketCreateRequest{
				FilmId:  1,
				PlaceId: 200,
			})

			//assert
			assert.Nil(t, resp)
			assert.EqualError(t, err, "rpc error: code = PermissionDenied desc = Header [token] is required")
		})

		t.Run("bad header token", func(t *testing.T) {
			//arrange
			Db.SetUp(t)
			defer Db.TearDown()

			ctx := context.Background()
			ctx = metadata.AppendToOutgoingContext(ctx, "Token", "any token")
			//в инит
			//act
			resp, err := ServerClient.TicketCreate(ctx, &server.TicketCreateRequest{
				FilmId:  1,
				PlaceId: 200,
			})

			//assert
			assert.Nil(t, resp)
			assert.EqualError(t, err, "rpc error: code = PermissionDenied desc = Header [Token] is invalid. See 'auth' method. Invalid Token. User not found by token")
		})

		t.Run("film not found", func(t *testing.T) {
			//arrange
			Db.SetUp(t)
			defer Db.TearDown()

			ctx := context.Background()

			// insert user
			var args []interface{}
			query := postgres.InsertUser
			user := internalModels.User{
				Id:    uint64(1),
				Name:  "UserName",
				Token: "UserToken",
			}
			args = append(args, user.Name)
			args = append(args, user.Token)
			Db.Pool.QueryRow(ctx, query, args...)

			ctx = metadata.AppendToOutgoingContext(ctx, "Token", user.Token)

			//act
			resp, err := ServerClient.TicketCreate(ctx, &server.TicketCreateRequest{
				FilmId:  1,
				PlaceId: 200,
			})

			//assert
			assert.Nil(t, resp)
			assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = Repository.CreateTicket.SelectFilm: film not found: no rows in result set")
		})

		t.Run("film room not found", func(t *testing.T) {
			//arrange
			Db.SetUp(t)
			defer Db.TearDown()

			ctx := context.Background()

			// insert user
			query := postgres.InsertUser
			user := internalModels.User{
				Id:    uint64(1),
				Name:  "UserName",
				Token: "UserToken",
			}
			var args []interface{}
			args = append(args, user.Name)
			args = append(args, user.Token)
			Db.Pool.QueryRow(ctx, query, args...)

			ctx = metadata.AppendToOutgoingContext(ctx, "Token", user.Token)

			// insert film
			query = postgres.InsertFilm
			film := internalModels.Film{
				Id:   uint64(1),
				Name: "Some film name",
			}
			args = nil
			args = append(args, film.Name)
			Db.Pool.QueryRow(ctx, query, args...)

			//act
			resp, err := ServerClient.TicketCreate(ctx, &server.TicketCreateRequest{
				FilmId:  1,
				PlaceId: 200,
			})

			//assert
			assert.Nil(t, resp)
			assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = Repository.CreateTicket.SelectFilmRoom: film_room not found: no rows in result set")
		})

		t.Run("place not exist in room", func(t *testing.T) {
			//arrange
			Db.SetUp(t)
			defer Db.TearDown()

			ctx := context.Background()
			countPlaces := 10
			insertPlaceId := uint64(countPlaces + 1)

			// insert user
			query := postgres.InsertUser
			user := internalModels.User{
				Id:    uint64(1),
				Name:  "UserName",
				Token: "UserToken",
			}
			var args []interface{}
			args = append(args, user.Name)
			args = append(args, user.Token)
			Db.Pool.QueryRow(ctx, query, args...)

			ctx = metadata.AppendToOutgoingContext(ctx, "Token", user.Token)

			// insert film
			query = postgres.InsertFilm
			film := internalModels.Film{
				Id:   uint64(1),
				Name: "Some film name",
			}
			args = nil
			args = append(args, film.Name)
			Db.Pool.QueryRow(ctx, query, args...)

			// insert room
			query = postgres.InsertRoom
			room := internalModels.Room{
				Id:     uint64(1),
				Places: []internalModels.Place{},
			}

			args = nil
			args = append(args, countPlaces)
			Db.Pool.QueryRow(ctx, query, args...)

			// insert film_room
			query = postgres.InsertFilmRoom
			args = nil
			args = append(args, film.Id)
			args = append(args, room.Id)
			Db.Pool.QueryRow(ctx, query, args...)

			//act
			resp, err := ServerClient.TicketCreate(ctx, &server.TicketCreateRequest{
				FilmId:  1,
				PlaceId: insertPlaceId,
			})

			//assert
			assert.Nil(t, resp)
			assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = Repository.CreateTicket.SelectFilmRoom: place not exist in room")
		})
	})

	t.Run("success", func(t *testing.T) {
		//arrange
		Db.SetUp(t)
		defer Db.TearDown()

		ctx := context.Background()
		countPlaces := 10
		insertPlaceId := uint64(countPlaces - 1)

		// insert user
		query := postgres.InsertUser
		user := internalModels.User{
			Id:    uint64(1),
			Name:  "UserName",
			Token: "UserToken",
		}
		var args []interface{}
		args = append(args, user.Name)
		args = append(args, user.Token)
		Db.Pool.QueryRow(ctx, query, args...)

		ctx = metadata.AppendToOutgoingContext(ctx, "Token", user.Token)

		// insert film
		query = postgres.InsertFilm
		film := internalModels.Film{
			Id:   uint64(1),
			Name: "Some film name",
		}
		args = nil
		args = append(args, film.Name)
		Db.Pool.QueryRow(ctx, query, args...)

		// insert room
		query = postgres.InsertRoom
		room := internalModels.Room{
			Id:     uint64(1),
			Places: []internalModels.Place{},
		}

		args = nil
		args = append(args, countPlaces)
		Db.Pool.QueryRow(ctx, query, args...)

		// insert film_room
		query = postgres.InsertFilmRoom
		args = nil
		args = append(args, film.Id)
		args = append(args, room.Id)
		Db.Pool.QueryRow(ctx, query, args...)

		apiResp := &server.TicketCreateResponse{
			Ticket: &server.Ticket{
				Id:      uint64(1),
				FilmId:  uint64(1),
				RoomId:  uint64(1),
				PlaceId: insertPlaceId,
			},
		}

		//act
		resp, err := ServerClient.TicketCreate(ctx, &server.TicketCreateRequest{
			FilmId:  1,
			PlaceId: insertPlaceId,
		})

		//assert
		assert.Nil(t, err)
		assert.Equal(t, resp.Ticket, apiResp.Ticket)
	})

}
