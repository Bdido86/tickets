package storage

import (
	"github.com/pkg/errors"
	"math/rand"
)

var dataUsers map[uint]*User
var dataRooms map[uint]*Room

func init() {
	dataUsers = make(map[uint]*User)
	dataRooms = make(map[uint]*Room)

	for id, _ := range films {
		places := make(map[uint]*Place)
		maxSeat := 2 + rand.Intn(6)
		for i := 1; i < maxSeat; i++ {
			place := Place{
				number: uint(i),
			}
			places[uint(i)] = &place
		}
		room := Room{
			number: id,
			places: places,
		}

		dataRooms[id] = &room
	}
}

func GetRoom(filmId uint) (*Room, error) {
	room, ok := dataRooms[filmId]
	if !ok {
		return nil, errors.New("Bad film ID")
	}

	return room, nil
}

func StartUser(id uint, name string) {
	if _, ok := dataUsers[id]; ok {
		return
	}

	user := newUser(id, name)
	dataUsers[user.id] = user
}

func GetFilms() map[uint]string {
	return films
}
