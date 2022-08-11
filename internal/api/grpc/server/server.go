package grpc

import (
	"context"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/repository"
	pb "gitlab.ozon.dev/Bdido86/movie-tickets/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedCinemaServer
	Deps
}

type Deps struct {
	CinemaRepository repository.Cinema
}

func NewServer(d Deps) pb.CinemaServer {
	return &server{
		Deps: d,
	}
}

func (s *server) UserAuth(ctx context.Context, in *pb.UserAuthRequest) (*pb.UserAuthResponse, error) {
	userName := in.GetName()
	if len(userName) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Field: [name] is required")
	}

	user, err := s.CinemaRepository.AuthUser(ctx, userName)
	if err != nil {
		return &pb.UserAuthResponse{}, errors.Wrap(err, "Error AuthUser")
	}

	return &pb.UserAuthResponse{
		Token: user.Token,
	}, nil
}

func (s *server) Films(in *pb.FilmsRequest, stream pb.Cinema_FilmsServer) error {
	limit64 := in.GetLimit()
	offset64 := in.GetOffset()
	desc := in.GetDesc()

	err := s.CinemaRepository.GetFilms(
		stream.Context(),
		limit64,
		offset64,
		desc,
		func(film *pb.Film) error {
			res := &pb.FilmsResponse{Film: film}
			err := stream.Send(res)
			if err != nil {
				return err
			}
			return nil
		},
	)
	if err != nil {
		return errors.Wrap(err, "Error Films")
	}

	return nil
}

func (s *server) FilmRoom(ctx context.Context, in *pb.FilmRoomRequest) (*pb.FilmRoomResponse, error) {
	film64 := in.GetFilmId()
	filmId := uint(film64)

	filmRoom, err := s.CinemaRepository.GetFilmRoom(ctx, filmId, getCurrentUserId(ctx))
	if err != nil {
		return &pb.FilmRoomResponse{}, errors.Wrap(err, "Error FilmRoom")
	}

	placesResponse := make([]*pb.FilmRoomResponse_Place, 0, len(filmRoom.Room.Places))
	for _, place := range filmRoom.Room.Places {
		placesResponse = append(placesResponse, &pb.FilmRoomResponse_Place{
			Id:     place.Id,
			IsFree: place.IsFree,
			IsMy:   place.IsMy,
		})
	}

	FilmResponse := &pb.Film{
		Id:   filmRoom.Film.Id,
		Name: filmRoom.Film.Name,
	}
	RoomResponse := &pb.FilmRoomResponse_Room{
		Id:     filmRoom.Room.Id,
		Places: placesResponse,
	}

	return &pb.FilmRoomResponse{
		Film: FilmResponse,
		Room: RoomResponse,
	}, nil
}

func (s *server) TicketCreate(ctx context.Context, in *pb.TicketCreateRequest) (*pb.TicketCreateResponse, error) {
	film64 := in.GetFilmId()
	place64 := in.GetPlaceId()

	filmId := uint(film64)
	placeId := uint(place64)

	ticket, err := s.CinemaRepository.CreateTicket(ctx, filmId, placeId, getCurrentUserId(ctx))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.TicketCreateResponse{
		Ticket: &pb.Ticket{
			Id:      ticket.Id,
			FilmId:  ticket.FilmId,
			RoomId:  ticket.RoomId,
			PlaceId: ticket.Place,
		},
	}, nil
}

func (s *server) TicketDelete(ctx context.Context, in *pb.TicketDeleteRequest) (*pb.TicketDeleteResponse, error) {
	ticket64 := in.GetTicketId()
	ticketId := uint(ticket64)

	err := s.CinemaRepository.DeleteTicket(ctx, ticketId, getCurrentUserId(ctx))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.TicketDeleteResponse{}, nil
}

func (s *server) MyTickets(ctx context.Context, _ *pb.MyTicketsRequest) (*pb.MyTicketsResponse, error) {
	tickets, err := s.CinemaRepository.GetMyTickets(ctx, getCurrentUserId(ctx))
	if err != nil {
		return &pb.MyTicketsResponse{}, errors.Wrap(err, "Error MyTickets")
	}

	ticketsResponse := make([]*pb.Ticket, 0, len(tickets))
	for _, ticket := range tickets {
		ticketsResponse = append(ticketsResponse, &pb.Ticket{
			Id:      ticket.Id,
			FilmId:  ticket.FilmId,
			RoomId:  ticket.RoomId,
			PlaceId: ticket.Place,
		})
	}

	return &pb.MyTicketsResponse{
		Tickets: ticketsResponse,
	}, nil
}

func getCurrentUserId(ctx context.Context) uint {
	return ctx.Value("userId").(uint)
}
