package grpc

import (
	"context"
	"github.com/golang/mock/gomock"
	mock_server "gitlab.ozon.dev/Bdido86/movie-tickets/internal/api/grpc/server/mocks"
	mock_broker "gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/broker/mocks"
	"testing"
)

type serverFixture struct {
	ctx         context.Context
	ctxWithUser context.Context
	mockServer  *mock_server.MockCinemaBackendClient
	mockBroker  *mock_broker.MockBroker
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
	f.mockBroker = mock_broker.NewMockBroker(mockCtrl)

	deps := Deps{
		Server: f.mockServer,
		Broker: f.mockBroker,
	}
	f.server = NewServer(deps)
	return f
}
