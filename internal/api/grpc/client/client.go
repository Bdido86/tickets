package grpc

import (
	"context"
	pb "gitlab.ozon.dev/Bdido86/movie-tickets/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"strings"
)

type server struct {
	pb.UnimplementedCinemaServer
	Deps
}

type Deps struct {
	Client pb.CinemaClient
}

func NewServer(d Deps) pb.CinemaServer {
	return &server{
		Deps: d,
	}
}

func (s *server) UserAuth(ctx context.Context, in *pb.UserAuthRequest) (*pb.UserAuthResponse, error) {
	userName := strings.Trim(in.GetName(), " ")
	if len(userName) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Field: [name] is required")
	}
	if len(userName) < 2 || len(userName) > 40 {
		return nil, status.Error(codes.InvalidArgument, "Field: [name] must be between 2-40 chars!")
	}

	return s.Client.UserAuth(ctx, in)
}

func (s *server) Films(in *pb.FilmsRequest, stream pb.Cinema_FilmsServer) error {
	limit64 := in.GetLimit()
	offset64 := in.GetOffset()
	if limit64 > 100 {
		return status.Error(codes.InvalidArgument, "Field: [limit] is too big. Maximum 100")
	}
	if offset64 > 100 {
		return status.Error(codes.InvalidArgument, "Field: [offset] is too big. Maximum 100")
	}

	ctx := stream.Context()
	ctx = prepareContext(ctx)
	streamFilm, err := s.Client.Films(ctx, in)
	if err != nil {
		return status.Error(codes.Unavailable, "Cannot get films stream")
	}
	for {
		res, err := streamFilm.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return status.Error(codes.Unavailable, "Cannot receive response stream")
		}

		film := res.GetFilm()
		resFilm := &pb.FilmsResponse{Film: film}
		err = stream.Send(resFilm)
		if err != nil {
			return status.Error(codes.Unavailable, "Cannot receive response stream film")
		}
	}

	return err
}

func (s *server) FilmRoom(ctx context.Context, in *pb.FilmRoomRequest) (*pb.FilmRoomResponse, error) {
	film64 := in.GetFilmId()
	if film64 == 0 {
		return nil, status.Error(codes.InvalidArgument, "Field: [filmId] is required and > 0")
	}
	if film64 > 20 {
		return nil, status.Error(codes.InvalidArgument, "Field: [filmId] is too big. Maximum 20")
	}

	ctx = prepareContext(ctx)
	return s.Client.FilmRoom(ctx, in)
}

func (s *server) TicketCreate(ctx context.Context, in *pb.TicketCreateRequest) (*pb.TicketCreateResponse, error) {
	film64 := in.GetFilmId()
	place64 := in.GetPlaceId()
	if film64 == 0 {
		return nil, status.Error(codes.InvalidArgument, "Field: [filmId] is required and > 0")
	}
	if film64 > 20 {
		return nil, status.Error(codes.InvalidArgument, "Field: [filmId] is too big. Maximum 20")
	}
	if place64 == 0 {
		return nil, status.Error(codes.InvalidArgument, "Field: [placeId] is required and > 0")
	}
	if place64 > 50 {
		return nil, status.Error(codes.InvalidArgument, "Field: [placeId] is too big. Maximum 50")
	}

	ctx = prepareContext(ctx)
	return s.Client.TicketCreate(ctx, in)
}

func (s *server) TicketDelete(ctx context.Context, in *pb.TicketDeleteRequest) (*pb.TicketDeleteResponse, error) {
	ticket64 := in.GetTicketId()
	if ticket64 == 0 {
		return nil, status.Error(codes.InvalidArgument, "Field: [ticketId] is required and > 0")
	}
	if ticket64 > 500 {
		return nil, status.Error(codes.InvalidArgument, "Field: [placeId] is too big. Maximum 500")
	}

	ctx = prepareContext(ctx)
	return s.Client.TicketDelete(ctx, in)
}

func (s *server) MyTickets(ctx context.Context, in *pb.MyTicketsRequest) (*pb.MyTicketsResponse, error) {
	ctx = prepareContext(ctx)
	return s.Client.MyTickets(ctx, in)
}

func prepareContext(ctx context.Context) context.Context {
	metaData, _ := metadata.FromIncomingContext(ctx)
	return metadata.NewOutgoingContext(ctx, metaData)
}
