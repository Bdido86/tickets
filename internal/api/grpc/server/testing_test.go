package grpc

import (
	"context"
	"github.com/golang/mock/gomock"
	mock_repository "gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/repository/mocks"
	"testing"
)

type cinemaFixture struct {
	ctx         context.Context
	ctxWithUser context.Context
	mockRepo    *mock_repository.MockCinema
	service     *server
}

func serverSetUp(t *testing.T) cinemaFixture {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	f := cinemaFixture{}
	f.ctx = context.Background()
	f.ctxWithUser = context.WithValue(context.Background(), "userId", uint(1))
	f.mockRepo = mock_repository.NewMockCinema(mockCtrl)

	deps := Deps{
		CinemaRepository: f.mockRepo,
	}
	f.service = NewServer(deps)
	return f
}
