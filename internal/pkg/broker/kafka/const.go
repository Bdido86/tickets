package kafka

const (
	topicTicketCreate = "TicketCreate"
	topicTicketDelete = "TicketDelete"
)

var (
	brokers = []string{"localhost:9095", "localhost:9096", "localhost:9097"}
)

type createTicketStruct struct {
	FilmId  uint   `json:"film_id"`
	PlaceId uint   `json:"place_id"`
	Token   string `json:"token"`
}

type deleteTicketStruct struct {
	Id    uint   `json:"id"`
	Token string `json:"token"`
}
