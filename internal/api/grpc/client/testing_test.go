package grpc

import (
	"context"
	"github.com/golang/mock/gomock"
	mock_server "gitlab.ozon.dev/Bdido86/movie-tickets/internal/api/grpc/server/mocks"
	"testing"
)

type serverFixture struct {
	ctx         context.Context
	ctxWithUser context.Context
	mockServer  *mock_server.MockCinemaBackendClient
	server      *server
}

func clientSetUp(t *testing.T) serverFixture {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	f := serverFixture{}
	f.ctx = context.Background()
	f.ctxWithUser = context.WithValue(context.Background(), "userId", uint(1))
	f.mockServer = mock_server.NewMockCinemaBackendClient(mockCtrl)

	deps := Deps{
		Server: f.mockServer,
	}
	f.server = NewServer(deps)
	return f
}
