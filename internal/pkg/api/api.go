package api

import (
	"context"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/storage"
	pb "gitlab.ozon.dev/Bdido86/movie-tickets/pkg/api"
)

func New() pb.CinemaServer {
	return &implementation{}
}

type implementation struct {
	pb.UnimplementedCinemaServer
}

func (i *implementation) UserAuth(_ context.Context, in *pb.UserAuthRequest) (*pb.UserAuthResponse, error) {
	user := in.GetName()
	if len(user) == 0 {
		return nil, errors.New("Field: [name] is required")
	}

	token := storage.AuthUser(user)
	resp := &pb.UserAuthResponse{
		Token: token,
	}

	return resp, nil
}
