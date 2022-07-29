package storage

var lastTicketId = uint(0)

var films = map[uint]string{
	1: "Зеленая миля",
	2: "Бойцовский клуб",
	3: "Джентельмены удачи",
	4: "Титаник",
	5: "Вверх (мультфильм)",
	6: "Один дома 2",
	7: "Курсы GO в ozon{ech (обучение)",
}

type Room struct {
	number uint
	places map[uint]*Place
}

type Place struct {
	number uint
	userId uint
}

type Ticket struct {
	id      uint
	roomId  uint
	placeId uint
	filmId  uint
}

func newTicket(filmId, roomId, placeId uint) *Ticket {
	lastTicketId++

	return &Ticket{
		id:      lastTicketId,
		filmId:  filmId,
		roomId:  roomId,
		placeId: placeId,
	}
}

func (r Room) GetNumber() uint {
	return r.number
}

func (r Room) GetPlaces() map[uint]*Place {
	return r.places
}

func (p Place) IsFree() bool {
	return p.userId == uint(0)
}

func (p Place) GetNumber() uint {
	return p.number
}

func (p Place) GetUserId() uint {
	return p.userId
}

func (t Ticket) GetId() uint {
	return t.id
}

func (t Ticket) GetPlaceId() uint {
	return t.placeId
}

func (t Ticket) GetRoomId() uint {
	return t.roomId
}

func (t Ticket) GetFilmId() uint {
	return t.filmId
}
