package grpc

import (
	"context"
	"errors"
	pbServerApi "gitlab.ozon.dev/Bdido86/movie-tickets/pkg/api/server"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/models"
)

func TestUserAuth(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// arrange
		name := "test-user"
		token := "JDJhJDEwJGcwSXRIYXlGZ2pKSGNUVnRxS0U1eE84aTFhTzdrQjZ2VWw3QmI1bEFOM1RUN2E4YmhqUmRL"
		modelUser := models.User{
			Id:    1,
			Name:  name,
			Token: token,
		}

		inApi := &pbServerApi.UserAuthRequest{
			Name: name,
		}
		respApi := &pbServerApi.UserAuthResponse{
			Token: token,
		}

		f := serverSetUp(t)
		f.mockRepo.EXPECT().AuthUser(f.ctx, name).Return(modelUser, nil).Times(1)

		// act
		resp, err := f.service.UserAuth(f.ctx, inApi)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, resp, respApi)
	})

	t.Run("error", func(t *testing.T) {
		t.Run("empty name", func(t *testing.T) {
			// arrange
			f := serverSetUp(t)
			inApi := &pbServerApi.UserAuthRequest{
				Name: "",
			}
			respErr := "rpc error: code = InvalidArgument desc = Field: [name] is required"

			// act
			_, err := f.service.UserAuth(context.Background(), inApi)

			// assert
			assert.EqualError(t, err, respErr)
		})

		t.Run("some error", func(t *testing.T) {
			// arrange
			name := "test-error"
			modelUser := models.User{}
			repoError := errors.New("some error")

			inApi := &pbServerApi.UserAuthRequest{
				Name: name,
			}
			respApi := &pbServerApi.UserAuthResponse{}
			respErr := "Error AuthUser: some error"

			f := serverSetUp(t)
			f.mockRepo.EXPECT().AuthUser(f.ctx, name).Return(modelUser, repoError).Times(1)

			// act
			resp, err := f.service.UserAuth(f.ctx, inApi)

			// assert
			assert.EqualError(t, err, respErr)
			assert.Equal(t, resp, respApi)
		})
	})
}

func TestFilmRoom(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// arrange
		f := serverSetUp(t)

		filmId := uint(1)
		userId := f.ctxWithUser.Value("userId").(uint)
		modelFilmRoom := models.FilmRoom{
			Film: models.Film{
				Id:   1,
				Name: "Фильм",
			},
			Room: models.Room{
				Id: 1,
				Places: []models.Place{
					{
						Id:     1,
						IsMy:   false,
						IsFree: true,
					},
					{
						Id:     2,
						IsMy:   true,
						IsFree: false,
					},
				},
			},
		}
		f.mockRepo.EXPECT().GetFilmRoom(f.ctxWithUser, filmId, userId).Return(modelFilmRoom, nil).Times(1)

		inApi := &pbServerApi.FilmRoomRequest{
			FilmId: uint64(filmId),
		}
		respApi := &pbServerApi.FilmRoomResponse{
			Film: &pbServerApi.Film{
				Id:   modelFilmRoom.Film.Id,
				Name: modelFilmRoom.Film.Name,
			},
			Room: &pbServerApi.FilmRoomResponse_Room{
				Id: modelFilmRoom.Room.Id,
				Places: []*pbServerApi.FilmRoomResponse_Place{
					{
						Id:     modelFilmRoom.Room.Places[0].Id,
						IsMy:   modelFilmRoom.Room.Places[0].IsMy,
						IsFree: modelFilmRoom.Room.Places[0].IsFree,
					},
					{
						Id:     modelFilmRoom.Room.Places[1].Id,
						IsMy:   modelFilmRoom.Room.Places[1].IsMy,
						IsFree: modelFilmRoom.Room.Places[1].IsFree,
					},
				},
			},
		}

		// act
		resp, err := f.service.FilmRoom(f.ctxWithUser, inApi)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, resp, respApi)
	})

	t.Run("error", func(t *testing.T) {
		// arrange
		f := serverSetUp(t)

		filmId := uint(1)
		userId := f.ctxWithUser.Value("userId").(uint)
		modelFilmRoom := models.FilmRoom{
			Film: models.Film{
				Id:   1,
				Name: "Фильм",
			},
			Room: models.Room{
				Id: 1,
				Places: []models.Place{
					{
						Id:     1,
						IsMy:   false,
						IsFree: true,
					},
				},
			},
		}
		repoError := errors.New("some error")
		f.mockRepo.EXPECT().GetFilmRoom(f.ctxWithUser, filmId, userId).Return(modelFilmRoom, repoError).Times(1)

		respErr := "Error FilmRoom: some error"
		inApi := &pbServerApi.FilmRoomRequest{
			FilmId: uint64(filmId),
		}
		respApi := &pbServerApi.FilmRoomResponse{}

		// act
		resp, err := f.service.FilmRoom(f.ctxWithUser, inApi)

		// assert
		assert.EqualError(t, err, respErr)
		assert.Equal(t, resp, respApi)
	})
}

func TestTicketCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// arrange
		f := serverSetUp(t)

		filmId := uint(1)
		placeId := uint(5)
		userId := f.ctxWithUser.Value("userId").(uint)

		modelTicket := models.Ticket{
			Id:     uint64(1),
			UserId: uint64(userId),
			FilmId: uint64(filmId),
			RoomId: uint64(placeId),
			Place:  uint64(placeId),
		}

		inApi := &pbServerApi.TicketCreateRequest{
			PlaceId: uint64(placeId),
			FilmId:  uint64(filmId),
		}
		respApi := &pbServerApi.TicketCreateResponse{
			Ticket: &pbServerApi.Ticket{
				Id:      modelTicket.Id,
				FilmId:  modelTicket.FilmId,
				RoomId:  modelTicket.RoomId,
				PlaceId: modelTicket.Place,
			},
		}
		f.mockRepo.EXPECT().CreateTicket(f.ctxWithUser, filmId, placeId, userId).Return(modelTicket, nil).Times(1)

		// act
		resp, err := f.service.TicketCreate(f.ctxWithUser, inApi)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, resp, respApi)
	})

	t.Run("error", func(t *testing.T) {
		// arrange
		f := serverSetUp(t)

		var modelTicket models.Ticket
		filmId := uint(1)
		placeId := uint(5)
		userId := f.ctxWithUser.Value("userId").(uint)

		repoError := errors.New("some error")
		inApi := &pbServerApi.TicketCreateRequest{
			PlaceId: uint64(placeId),
			FilmId:  uint64(filmId),
		}

		f.mockRepo.EXPECT().CreateTicket(f.ctxWithUser, filmId, placeId, userId).Return(modelTicket, repoError).Times(1)
		respErr := "rpc error: code = InvalidArgument desc = some error"

		// act
		resp, err := f.service.TicketCreate(f.ctxWithUser, inApi)

		// assert
		assert.EqualError(t, err, respErr)
		assert.Nil(t, resp)
	})
}

func TestTicketDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// arrange
		ticketId := uint(1)
		inApi := &pbServerApi.TicketDeleteRequest{
			TicketId: uint64(ticketId),
		}

		respApi := &pbServerApi.TicketDeleteResponse{}
		f := serverSetUp(t)
		f.mockRepo.EXPECT().DeleteTicket(f.ctxWithUser, ticketId, f.ctxWithUser.Value("userId")).Return(nil).Times(1)

		// act
		resp, err := f.service.TicketDelete(f.ctxWithUser, inApi)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, resp, respApi)
	})

	t.Run("error", func(t *testing.T) {
		// arrange
		ticketId := uint(1)
		inApi := &pbServerApi.TicketDeleteRequest{
			TicketId: uint64(ticketId),
		}

		repoError := errors.New("some error")
		respErr := "rpc error: code = InvalidArgument desc = some error"

		f := serverSetUp(t)
		f.mockRepo.EXPECT().DeleteTicket(f.ctxWithUser, ticketId, f.ctxWithUser.Value("userId")).Return(repoError).Times(1)

		// act
		resp, err := f.service.TicketDelete(f.ctxWithUser, inApi)

		// assert
		assert.EqualError(t, err, respErr)
		assert.Nil(t, resp)
	})
}

func TestMyTickets(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Run("list", func(t *testing.T) {
			// arrange
			modelTickets := []models.Ticket{
				{
					Id:     1,
					UserId: 1,
					FilmId: 1,
					RoomId: 1,
					Place:  1,
				},
				{
					Id:     2,
					UserId: 1,
					FilmId: 1,
					RoomId: 1,
					Place:  2,
				},
				{
					Id:     3,
					UserId: 1,
					FilmId: 2,
					RoomId: 2,
					Place:  5,
				},
			}
			respApi := &pbServerApi.MyTicketsResponse{
				Tickets: []*pbServerApi.Ticket{
					{
						Id:      1,
						FilmId:  1,
						RoomId:  1,
						PlaceId: 1,
					},
					{
						Id:      2,
						FilmId:  1,
						RoomId:  1,
						PlaceId: 2,
					},
					{
						Id:      3,
						FilmId:  2,
						RoomId:  2,
						PlaceId: 5,
					},
				},
			}

			f := serverSetUp(t)
			f.mockRepo.EXPECT().GetMyTickets(f.ctxWithUser, f.ctxWithUser.Value("userId")).Return(modelTickets, nil).Times(1)

			// act
			resp, err := f.service.MyTickets(f.ctxWithUser, &pbServerApi.MyTicketsRequest{})

			// assert
			assert.NoError(t, err)
			assert.Equal(t, resp, respApi)
		})

		t.Run("one", func(t *testing.T) {
			// arrange
			modelTickets := []models.Ticket{
				{
					Id:     10,
					UserId: 1,
					FilmId: 10,
					RoomId: 10,
					Place:  3,
				},
			}
			respApi := &pbServerApi.MyTicketsResponse{
				Tickets: []*pbServerApi.Ticket{
					{
						Id:      10,
						FilmId:  10,
						RoomId:  10,
						PlaceId: 3,
					},
				},
			}

			f := serverSetUp(t)
			f.mockRepo.EXPECT().GetMyTickets(f.ctxWithUser, f.ctxWithUser.Value("userId")).Return(modelTickets, nil).Times(1)

			// act
			resp, err := f.service.MyTickets(f.ctxWithUser, &pbServerApi.MyTicketsRequest{})

			// assert
			assert.NoError(t, err)
			assert.Equal(t, resp, respApi)
		})

		t.Run("empty", func(t *testing.T) {
			// arrange
			var modelTickets []models.Ticket
			respApi := &pbServerApi.MyTicketsResponse{
				Tickets: []*pbServerApi.Ticket{},
			}

			f := serverSetUp(t)
			f.mockRepo.EXPECT().GetMyTickets(f.ctxWithUser, f.ctxWithUser.Value("userId")).Return(modelTickets, nil).Times(1)

			// act
			resp, err := f.service.MyTickets(f.ctxWithUser, &pbServerApi.MyTicketsRequest{})

			// assert
			assert.NoError(t, err)
			assert.Equal(t, resp, respApi)
		})
	})

	t.Run("error", func(t *testing.T) {
		// arrange
		modelTickets := []models.Ticket{}
		repoError := errors.New("some error")
		respApi := &pbServerApi.MyTicketsResponse{}

		f := serverSetUp(t)
		f.mockRepo.EXPECT().GetMyTickets(f.ctxWithUser, f.ctxWithUser.Value("userId")).Return(modelTickets, repoError).Times(1)
		respErr := "Error MyTickets: some error"

		// act
		resp, err := f.service.MyTickets(f.ctxWithUser, &pbServerApi.MyTicketsRequest{})

		// assert
		assert.EqualError(t, err, respErr)
		assert.Equal(t, resp, respApi)

	})
}

func TestGetCurrentUserId(t *testing.T) {
	f := serverSetUp(t)

	// arrange
	tests := []struct {
		ctx  context.Context
		resp uint
	}{
		{
			ctx:  f.ctxWithUser,
			resp: 1,
		},
		{
			ctx:  context.WithValue(context.Background(), "userId", uint(11)),
			resp: 11,
		},
	}

	for _, testData := range tests {
		// act
		resp := getCurrentUserId(testData.ctx)

		// assert
		assert.Equal(t, testData.resp, resp)
	}
}
