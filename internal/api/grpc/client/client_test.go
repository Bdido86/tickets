package grpc

import (
	"context"
	"errors"
	pbClientApi "gitlab.ozon.dev/Bdido86/movie-tickets/pkg/api/client"
	pbServerApi "gitlab.ozon.dev/Bdido86/movie-tickets/pkg/api/server"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserAuth(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		t.Run("empty name", func(t *testing.T) {
			// arrange
			f := clientSetUp(t)

			inApi := &pbClientApi.UserAuthRequest{
				Name: "",
			}
			respErr := "rpc error: code = InvalidArgument desc = Field: [name] is required"

			// act
			_, err := f.server.UserAuth(context.Background(), inApi)

			// assert
			assert.EqualError(t, err, respErr)
		})

		t.Run("short name", func(t *testing.T) {
			// arrange
			f := clientSetUp(t)

			inApi := &pbClientApi.UserAuthRequest{
				Name: "x",
			}
			respErr := "rpc error: code = InvalidArgument desc = Field: [name] must be between 2-40 chars!"

			// act
			_, err := f.server.UserAuth(context.Background(), inApi)

			// assert
			assert.EqualError(t, err, respErr)
		})

		t.Run("long name", func(t *testing.T) {
			// arrange
			f := clientSetUp(t)

			inApi := &pbClientApi.UserAuthRequest{
				Name: "name1234567890123456789012345678901234567890end",
			}
			respErr := "rpc error: code = InvalidArgument desc = Field: [name] must be between 2-40 chars!"

			// act
			_, err := f.server.UserAuth(context.Background(), inApi)

			// assert
			assert.EqualError(t, err, respErr)
		})

		t.Run("server some error", func(t *testing.T) {
			// arrange
			f := clientSetUp(t)

			name := "name"
			inApiClient := &pbClientApi.UserAuthRequest{
				Name: name,
			}
			inServerApi := &pbServerApi.UserAuthRequest{
				Name: name,
			}

			respErr := "some error"
			serverError := errors.New(respErr)
			apiRespServer := &pbServerApi.UserAuthResponse{}
			f.mockServer.EXPECT().UserAuth(f.ctx, inServerApi).Return(apiRespServer, serverError).Times(1)

			// act
			resp, err := f.server.UserAuth(context.Background(), inApiClient)

			// assert
			assert.EqualError(t, err, respErr)
			assert.Nil(t, resp)
		})
	})

	t.Run("success", func(t *testing.T) {
		// arrange
		f := clientSetUp(t)

		name := "name"
		inApiClient := &pbClientApi.UserAuthRequest{
			Name: name,
		}
		inServerApi := &pbServerApi.UserAuthRequest{
			Name: name,
		}

		token := "someToken"
		apiRespServer := &pbServerApi.UserAuthResponse{
			Token: token,
		}
		apiRespClient := &pbClientApi.UserAuthResponse{
			Token: token,
		}
		f.mockServer.EXPECT().UserAuth(f.ctx, inServerApi).Return(apiRespServer, nil).Times(1)

		// act
		resp, err := f.server.UserAuth(context.Background(), inApiClient)

		// assert
		assert.Equal(t, resp, apiRespClient)
		assert.Nil(t, err)
	})

}
