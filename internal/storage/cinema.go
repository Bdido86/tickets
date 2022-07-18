package storage

var films = map[uint]string{
	1: "Зеленая миля",
	2: "Бойцовский клуб",
	3: "Джентельмены удачи",
	4: "Титаник",
	5: "Вверх (мультфильм)",
	6: "Курсы GO в ozon{ech (обучение)",
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
