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

func TestAuthUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		Db.SetUp(t)
		defer Db.TearDown()

		var (
			name      = "dido"
			existUser internalModels.User
		)

		t.Run("new user", func(t *testing.T) {
			//arrange

			// act
			res, err := Repository.AuthUser(context.Background(), name)

			// assert
			assert.NoError(t, err)
			assert.Equal(t, res, internalModels.User{
				Id:    uint64(1),
				Name:  "dido",
				Token: res.Token,
			})

			existUser = res
		})

		t.Run("existing user", func(t *testing.T) {
			//arrange

			// act
			res, err := Repository.AuthUser(context.Background(), name)

			// assert
			assert.NoError(t, err)
			assert.Equal(t, res, existUser)
		})
	})
}

func TestGetUserIdByToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		Db.SetUp(t)
		defer Db.TearDown()

		//arrange
		ctx := context.Background()

		apiResp := internalModels.User{
			Id:    uint64(1),
			Name:  "userName",
			Token: "userToken",
		}

		// insert user
		query := postgres.InsertUser
		var args []interface{}
		args = append(args, apiResp.Name)
		args = append(args, apiResp.Token)
		Db.Pool.QueryRow(ctx, query, args...)

		// act
		res, err := Repository.GetUserIdByToken(ctx, apiResp.Token)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, res, uint(apiResp.Id))
	})

	t.Run("fail", func(t *testing.T) {
		Db.SetUp(t)
		defer Db.TearDown()

		//arrange
		token := "someToken"
		resUserId := uint(0)

		// act
		res, err := Repository.GetUserIdByToken(context.Background(), token)

		// assert
		assert.EqualError(t, err, "Invalid Token. User not found by token")
		assert.Equal(t, res, resUserId)
	})
}
