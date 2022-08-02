package api

import (
	"context"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/repository"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/storage"
	pb "gitlab.ozon.dev/Bdido86/movie-tickets/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strconv"
	"strings"
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
		return &pb.UserAuthResponse{}, errors.Wrap(err, "error AuthUser")
	}

	return &pb.UserAuthResponse{
		Token: user.Token,
	}, nil
}

func (s *server) Films(ctx context.Context, _ *pb.FilmsRequest) (*pb.FilmsResponse, error) {
	films, err := s.CinemaRepository.GetFilms(ctx)
	if err != nil {
		return &pb.FilmsResponse{}, errors.Wrap(err, "error Films")
	}

	result := make([]*pb.Film, 0, len(films))
	for _, film := range films {
		result = append(result, &pb.Film{
			Id:   uint64(film.Id),
			Name: film.Name,
		})
	}
	return &pb.FilmsResponse{
		Films: result,
	}, nil
}

func (s *server) FilmRoom(ctx context.Context, in *pb.FilmRoomRequest) (*pb.FilmRoomResponse, error) {
	film64 := in.GetFilmId()
	filmId := uint(film64)

	filmRoom, err := s.CinemaRepository.GetFilmRoom(ctx, filmId, getCurrentUserId(ctx))
	if err != nil {
		return &pb.FilmRoomResponse{}, errors.Wrap(err, "error FilmRoom")
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

	films := storage.GetFilms()
	_, ok := films[filmId]
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "Field: [film_id] not found")
	}

	ticket, err := storage.BuyTicket(filmId, placeId, getUserIdFromToken(ctx))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.TicketCreateResponse{
		Ticket: &pb.Ticket{
			Id:      uint64(ticket.GetId()),
			FilmId:  uint64(ticket.GetFilmId()),
			RoomId:  uint64(ticket.GetRoomId()),
			PlaceId: uint64(ticket.GetPlaceId()),
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

func getUserIdFromToken(ctx context.Context) uint {
	metaData, _ := metadata.FromIncomingContext(ctx)

	tokens := metaData.Get("Token")
	words := strings.Split(tokens[0], "-")
	id, _ := strconv.Atoi(words[1])

	return uint(id)
}

func getCurrentUserId(ctx context.Context) uint {
	return ctx.Value("userId").(uint)
}
