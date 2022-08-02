package models

type Film struct {
	Id   uint64 `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type User struct {
	Id    uint64 `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Token string `json:"token" db:"token"`
}

type Place struct {
	Id     uint64 `json:"id"`
	IsMy   bool   `json:"isMy"`
	IsFree bool   `json:"isFree"`
}

type Room struct {
	Id     uint64  `json:"id"`
	Places []Place `json:"places"`
}

type FilmRoom struct {
	Film Film `json:"film"`
	Room Room `json:"room"`
}

type RoomDb struct {
	Id          uint64 `json:"id" db:"id"`
	CountPlaces uint64 `json:"count_places" db:"count_places"`
}

type Ticket struct {
	Id     uint64 `json:"id" db:"id"`
	UserId uint64 `json:"user_id" db:"user_id"`
	FilmId uint64 `json:"film_id" db:"film_id"`
	RoomId uint64 `json:"room_id" db:"room_id"`
	Place  uint64 `json:"place" db:"place"`
}
