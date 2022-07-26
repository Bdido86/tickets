package api

import (
	"context"
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/storage"
	pb "gitlab.ozon.dev/Bdido86/movie-tickets/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"sort"
	"strconv"
	"strings"
)

type server struct {
	pb.UnimplementedCinemaServer
}

func NewServer() pb.CinemaServer {
	return &server{}
}

func (s *server) UserAuth(_ context.Context, in *pb.UserAuthRequest) (*pb.UserAuthResponse, error) {
	user := in.GetName()
	if len(user) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Field: [name] is required")
	}

	token := storage.AuthUser(user)
	return &pb.UserAuthResponse{
		Token: token,
	}, nil
}

func (s *server) Films(_ context.Context, _ *pb.FilmsRequest) (*pb.FilmsResponse, error) {
	films := storage.GetFilms()

	keys := make([]int, 0, len(films))
	for id, _ := range films {
		keys = append(keys, int(id))
	}
	sort.Ints(keys)

	result := make([]*pb.Film, 0, len(films))
	for _, k := range keys {
		result = append(result, &pb.Film{
			Id:   uint64(k),
			Name: films[uint(k)],
		})
	}

	return &pb.FilmsResponse{
		Films: result,
	}, nil
}

func (s *server) FilmRoom(ctx context.Context, in *pb.FilmRoomRequest) (*pb.FilmRoomResponse, error) {
	film64 := in.GetFilmId()
	filmId := uint(film64)
	films := storage.GetFilms()
	film, ok := films[filmId]
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "Field: [film_id] not found")
	}

	room, _ := storage.GetRoom(filmId)

	placesResponse := make([]*pb.FilmRoomResponse_Place, 0, len(room.GetPlaces()))
	for _, place := range room.GetPlaces() {
		placesResponse = append(placesResponse, &pb.FilmRoomResponse_Place{
			Id:     uint64(place.GetNumber()),
			IsFree: place.GetUserId() > 0,
			IsMy:   place.GetUserId() == getUserIdFromToken(ctx),
		})
	}

	FilmResponse := &pb.Film{
		Id:   film64,
		Name: film,
	}

	RoomResponse := &pb.FilmRoomResponse_Room{
		Id:     uint64(room.GetNumber()),
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

	err := storage.DeleteTicket(getUserIdFromToken(ctx), ticketId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Field: [ticket_id] not found")
	}

	return &pb.TicketDeleteResponse{}, nil
}

func (s *server) MyTickets(ctx context.Context, _ *pb.MyTicketsRequest) (*pb.MyTicketsResponse, error) {
	tickets, _ := storage.GetTickets(getUserIdFromToken(ctx))

	ticketsResponse := make([]*pb.Ticket, 0, len(tickets))
	for _, ticket := range tickets {
		ticketsResponse = append(ticketsResponse, &pb.Ticket{
			Id:      uint64(ticket.GetId()),
			FilmId:  uint64(ticket.GetFilmId()),
			RoomId:  uint64(ticket.GetRoomId()),
			PlaceId: uint64(ticket.GetPlaceId()),
		})
	}

	return &pb.MyTicketsResponse{
		Tickets: ticketsResponse,
	}, nil
}

func IsValidToken(token string) bool {
	return storage.IsValidToken(token)
}

func getUserIdFromToken(ctx context.Context) uint {
	metaData, _ := metadata.FromIncomingContext(ctx)

	tokens := metaData.Get("Token")
	words := strings.Split(tokens[0], "-")
	id, _ := strconv.Atoi(words[1])

	return uint(id)
}