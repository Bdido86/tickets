package grpc

import (
	"context"
	"encoding/json"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/broker"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/logger"
	pbClient "gitlab.ozon.dev/Bdido86/movie-tickets/pkg/api/client"
	pbServer "gitlab.ozon.dev/Bdido86/movie-tickets/pkg/api/server"
	"go.opencensus.io/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"strings"
)

type server struct {
	pbClient.UnimplementedCinemaFrontendServer
	Deps
}

type Deps struct {
	Server pbServer.CinemaBackendClient
	Broker broker.Broker
	Logger logger.Logger
}

func NewServer(d Deps) *server {
	return &server{
		Deps: d,
	}
}

func (s *server) UserAuth(ctx context.Context, in *pbClient.UserAuthRequest) (*pbClient.UserAuthResponse, error) {
	ctx, span := trace.StartSpan(ctx, "grpc/client/UserAuth")
	defer span.End()

	ctx = appendSpan(ctx, span)

	userName := strings.Trim(in.GetName(), " ")
	if len(userName) == 0 {
		s.Logger.Info("field: [name] is required")
		return nil, status.Error(codes.InvalidArgument, "Field: [name] is required")
	}
	if len(userName) < 2 || len(userName) > 40 {
		s.Logger.Info("field: [name] must be between 2-40 chars!")
		return nil, status.Error(codes.InvalidArgument, "Field: [name] must be between 2-40 chars!")
	}

	inServer := &pbServer.UserAuthRequest{
		Name: userName,
	}
	serverUserAuthResponse, err := s.Server.UserAuth(ctx, inServer)
	if err != nil {
		s.Logger.Errorf("error server UserAuth %v", err)
		return nil, err
	}

	return &pbClient.UserAuthResponse{
		Token: serverUserAuthResponse.Token,
	}, nil
}

func (s *server) Films(in *pbClient.FilmsRequest, stream pbClient.CinemaFrontend_FilmsServer) error {
	ctx := stream.Context()

	ctx, span := trace.StartSpan(ctx, "grpc/client/Films")
	defer span.End()

	ctx = appendSpan(ctx, span)

	limit64 := in.GetLimit()
	offset64 := in.GetOffset()
	if limit64 > 100 {
		s.Logger.Info("Field: [limit] is too big. Maximum 100")
		return status.Error(codes.InvalidArgument, "Field: [limit] is too big. Maximum 100")
	}
	if offset64 > 100 {
		s.Logger.Info("Field: [offset] is too big. Maximum 100")
		return status.Error(codes.InvalidArgument, "Field: [offset] is too big. Maximum 100")
	}

	inServer := &pbServer.FilmsRequest{
		Limit:  in.Limit,
		Offset: in.Offset,
		Desc:   in.Desc,
	}

	serverStreamFilm, err := s.Server.Films(ctx, inServer)
	if err != nil {
		s.Logger.Errorf("cannot get films stream %v", err)
		return status.Error(codes.Unavailable, "Cannot get films stream")
	}

	for {
		res, err := serverStreamFilm.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			s.Logger.Errorf("cannot receive response stream %v", err)
			return status.Error(codes.Unavailable, "Cannot receive response stream")
		}

		serverFilm := res.GetFilm()
		film := &pbClient.Film{
			Id:   serverFilm.Id,
			Name: serverFilm.Name,
		}

		resFilm := &pbClient.FilmsResponse{
			Film: film,
		}
		err = stream.Send(resFilm)
		if err != nil {
			s.Logger.Errorf("cannot receive response stream film %v", err)
			return status.Error(codes.Unavailable, "Cannot receive response stream film")
		}
	}

	return nil
}

func (s *server) FilmRoom(ctx context.Context, in *pbClient.FilmRoomRequest) (*pbClient.FilmRoomResponse, error) {
	ctx, span := trace.StartSpan(ctx, "grpc/client/FilmRoom")
	defer span.End()

	ctx = appendSpan(ctx, span)

	film64 := in.GetFilmId()
	if film64 == 0 {
		s.Logger.Info("Field: [filmId] is required and > 0")
		return nil, status.Error(codes.InvalidArgument, "Field: [filmId] is required and > 0")
	}
	if film64 > 20 {
		s.Logger.Info("Field: [filmId] is too big. Maximum 20")
		return nil, status.Error(codes.InvalidArgument, "Field: [filmId] is too big. Maximum 20")
	}

	inServer := &pbServer.FilmRoomRequest{
		FilmId: in.FilmId,
	}

	serverFilmRoomResponse, err := s.Server.FilmRoom(ctx, inServer)
	if err != nil {
		s.Logger.Errorf("server film %v", err)
		return nil, err
	}

	film := &pbClient.Film{
		Id:   serverFilmRoomResponse.Film.Id,
		Name: serverFilmRoomResponse.Film.Name,
	}

	places := make([]*pbClient.FilmRoomResponse_Place, 0, len(serverFilmRoomResponse.Room.Places))
	for _, serverPlace := range serverFilmRoomResponse.Room.Places {
		places = append(places, &pbClient.FilmRoomResponse_Place{
			Id:     serverPlace.Id,
			IsMy:   serverPlace.IsMy,
			IsFree: serverPlace.IsFree,
		})
	}

	return &pbClient.FilmRoomResponse{
		Film: film,
		Room: &pbClient.FilmRoomResponse_Room{
			Id:     serverFilmRoomResponse.Room.Id,
			Places: places,
		},
	}, nil
}

func (s *server) TicketCreate(ctx context.Context, in *pbClient.TicketCreateRequest) (*pbClient.TicketCreateResponse, error) {
	ctx, span := trace.StartSpan(ctx, "grpc/client/TicketCreate")
	defer span.End()

	ctx = appendSpan(ctx, span)

	film64 := in.GetFilmId()
	place64 := in.GetPlaceId()
	if film64 == 0 {
		s.Logger.Info("Field: [filmId] is required and > 0")
		return nil, status.Error(codes.InvalidArgument, "Field: [filmId] is required and > 0")
	}
	if film64 > 20 {
		s.Logger.Info("Field: [filmId] is too big. Maximum 20")
		return nil, status.Error(codes.InvalidArgument, "Field: [filmId] is too big. Maximum 20")
	}
	if place64 == 0 {
		s.Logger.Info("Field: [placeId] is required and > 0")
		return nil, status.Error(codes.InvalidArgument, "Field: [placeId] is required and > 0")
	}
	if place64 > 50 {
		s.Logger.Info("Field: [placeId] is too big. Maximum 50")
		return nil, status.Error(codes.InvalidArgument, "Field: [placeId] is too big. Maximum 50")
	}

	inServer := &pbServer.TicketCreateRequest{
		FilmId:  in.FilmId,
		PlaceId: in.PlaceId,
	}

	serverTicketCreateResponse, err := s.Server.TicketCreate(ctx, inServer)
	if err != nil {
		s.Logger.Errorf("server ticket create %v", err)
		return nil, err
	}

	return &pbClient.TicketCreateResponse{
		Ticket: &pbClient.Ticket{
			Id:      serverTicketCreateResponse.Ticket.Id,
			FilmId:  serverTicketCreateResponse.Ticket.FilmId,
			RoomId:  serverTicketCreateResponse.Ticket.RoomId,
			PlaceId: serverTicketCreateResponse.Ticket.PlaceId,
		},
	}, nil
}

func (s *server) TicketDelete(ctx context.Context, in *pbClient.TicketDeleteRequest) (*pbClient.TicketDeleteResponse, error) {
	ctx, span := trace.StartSpan(ctx, "grpc/client/TicketDelete")
	defer span.End()

	ctx = appendSpan(ctx, span)

	ticket64 := in.GetTicketId()
	if ticket64 == 0 {
		s.Logger.Info("Field: [ticketId] is required and > 0")
		return nil, status.Error(codes.InvalidArgument, "Field: [ticketId] is required and > 0")
	}
	if ticket64 > 500 {
		s.Logger.Info("Field: [ticket64] is too big. Maximum 500")
		return nil, status.Error(codes.InvalidArgument, "Field: [ticket64] is too big. Maximum 500")
	}

	inServer := &pbServer.TicketDeleteRequest{
		TicketId: in.TicketId,
	}

	_, err := s.Server.TicketDelete(ctx, inServer)
	if err != nil {
		s.Logger.Errorf("server ticket delete %v", err)
		return nil, err
	}

	return &pbClient.TicketDeleteResponse{}, nil
}

func (s *server) MyTickets(ctx context.Context, _ *pbClient.MyTicketsRequest) (*pbClient.MyTicketsResponse, error) {
	ctx, span := trace.StartSpan(ctx, "grpc/client/MyTickets")
	defer span.End()

	ctx = appendSpan(ctx, span)

	inServer := &pbServer.MyTicketsRequest{}
	serverMyTicketsResponse, err := s.Server.MyTickets(ctx, inServer)
	if err != nil {
		s.Logger.Errorf("server my tickets %v", err)
		return nil, err
	}

	tickets := make([]*pbClient.Ticket, 0, len(serverMyTicketsResponse.Tickets))
	for _, serverTicket := range serverMyTicketsResponse.Tickets {
		tickets = append(tickets, &pbClient.Ticket{
			Id:      serverTicket.Id,
			FilmId:  serverTicket.FilmId,
			RoomId:  serverTicket.RoomId,
			PlaceId: serverTicket.PlaceId,
		})
	}
	return &pbClient.MyTicketsResponse{
		Tickets: tickets,
	}, nil
}

func (s *server) TicketDeleteAsync(ctx context.Context, in *pbClient.TicketDeleteRequestAsync) (*pbClient.TicketDeleteResponseAsync, error) {
	ctx, span := trace.StartSpan(ctx, "grpc/client/TicketDeleteAsync")
	defer span.End()

	ctx = appendSpan(ctx, span)

	ticket64 := in.GetTicketId()
	if ticket64 == 0 {
		s.Logger.Info("Field: [ticketId] is required and > 0")
		return nil, status.Error(codes.InvalidArgument, "Field: [ticketId] is required and > 0")
	}
	if ticket64 > 500 {
		s.Logger.Info("Field: [ticketId] is too big. Maximum 500")
		return nil, status.Error(codes.InvalidArgument, "Field: [ticketId] is too big. Maximum 500")
	}

	err := s.Broker.DeleteTicket(ctx, uint(ticket64))
	if err != nil {
		s.Logger.Errorf("broker ticket delete %v", err)
		return nil, err
	}

	return &pbClient.TicketDeleteResponseAsync{}, nil
}

func (s *server) TicketCreateAsync(ctx context.Context, in *pbClient.TicketCreateRequestAsync) (*pbClient.TicketCreateResponseAsync, error) {
	ctx, span := trace.StartSpan(ctx, "grpc/client/TicketCreateAsync")
	defer span.End()

	ctx = appendSpan(ctx, span)

	film64 := in.GetFilmId()
	place64 := in.GetPlaceId()
	if film64 == 0 {
		s.Logger.Info("Field: [filmId] is required and > 0")
		return nil, status.Error(codes.InvalidArgument, "Field: [filmId] is required and > 0")
	}
	if film64 > 20 {
		s.Logger.Info("Field: [filmId] is too big. Maximum 20")
		return nil, status.Error(codes.InvalidArgument, "Field: [filmId] is too big. Maximum 20")
	}
	if place64 == 0 {
		s.Logger.Info("Field: [placeId] is required and > 0")
		return nil, status.Error(codes.InvalidArgument, "Field: [placeId] is required and > 0")
	}
	if place64 > 50 {
		s.Logger.Info("Field: [placeId] is too big. Maximum 50")
		return nil, status.Error(codes.InvalidArgument, "Field: [placeId] is too big. Maximum 50")
	}

	err := s.Broker.CreateTicket(ctx, uint(film64), uint(place64))
	if err != nil {
		s.Logger.Errorf("broker ticket create %v", err)
		return nil, err
	}

	return &pbClient.TicketCreateResponseAsync{}, nil
}

func appendSpan(ctx context.Context, span *trace.Span) context.Context {
	spanContextJson, _ := json.Marshal(span.SpanContext())

	metaData, _ := metadata.FromIncomingContext(ctx)
	metaData.Set("X-Span-Context", string(spanContextJson))

	return metadata.NewOutgoingContext(ctx, metaData)
}
