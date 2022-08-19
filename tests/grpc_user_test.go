//go:build integration
// +build integration

package tests

import (
	"context"
	"gitlab.ozon.dev/Bdido86/movie-tickets/pkg/api/server"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserAuth(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		t.Run("empty name", func(t *testing.T) {
			//arrange

			//в инит
			//act
			resp, err := ServerClient.UserAuth(context.Background(), &server.UserAuthRequest{
				Name: "",
			})

			//assert
			assert.Nil(t, resp)
			assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = Field: [name] is required")
		})
	})

	t.Run("success", func(t *testing.T) {
		Db.SetUp(t)
		defer Db.TearDown()

		var token string
		nameQa := "SomeNameQA"

		t.Run("new name", func(t *testing.T) {
			//arrange
			name := nameQa

			//act
			resp, err := ServerClient.UserAuth(context.Background(), &server.UserAuthRequest{
				Name: name,
			})
			token = resp.Token

			//assert
			assert.Nil(t, err)
			assert.NotEmpty(t, resp.Token)
		})

		t.Run("existing name", func(t *testing.T) {
			//arrange
			name := nameQa
			apiResp := &server.UserAuthResponse{
				Token: token,
			}

			//act
			resp, err := ServerClient.UserAuth(context.Background(), &server.UserAuthRequest{
				Name: name,
			})

			//assert
			assert.Nil(t, err)
			assert.Equal(t, resp.Token, apiResp.Token)
		})
	})
}
