package api

import (
	"context"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/storage"
	pb "gitlab.ozon.dev/Bdido86/movie-tickets/pkg/api"
	"sort"
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

func (i *implementation) Films(_ context.Context, _ *pb.FilmsRequest) (*pb.FilmsResponse, error) {
	films := storage.GetFilms()

	keys := make([]int, 0, len(films))
	for id, _ := range films {
		keys = append(keys, int(id))
	}
	sort.Ints(keys)

	result := make([]*pb.FilmsResponse_Film, 0, len(films))
	for _, k := range keys {
		result = append(result, &pb.FilmsResponse_Film{
			Id:   uint64(k),
			Name: films[uint(k)],
		})
	}

	return &pb.FilmsResponse{
		Films: result,
	}, nil
}

func IsValidToken(token string) bool {
	return storage.IsValidToken(token)
}
