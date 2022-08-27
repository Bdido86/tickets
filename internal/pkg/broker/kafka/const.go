package kafka

const (
	topicTicketCreate = "TicketCreate"
	topicTicketDelete = "TicketDelete"

	token = "Token"
)

var (
	brokers = []string{"localhost:9095", "localhost:9096", "localhost:9097"}
)

type createTicketStruct struct {
	FilmId  uint `json:"film_id"`
	PlaceId uint `json:"place_id"`
}

type deleteTicketStruct struct {
	Id uint `json:"id"`
}
