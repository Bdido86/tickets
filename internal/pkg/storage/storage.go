package storage

import (
	"github.com/pkg/errors"
	"math/rand"
	"sync"
)

const poolSize = 10

var poolChannel chan struct{}
var dataUsers map[uint]*User
var dataRooms map[uint]*Room
var mu sync.RWMutex

func init() {
	mu = sync.RWMutex{}
	poolChannel = make(chan struct{}, poolSize)
	dataUsers = make(map[uint]*User)
	dataRooms = make(map[uint]*Room)

	for id, _ := range films {
		places := make(map[uint]*Place)
		maxSeat := 2 + rand.Intn(6)
		for i := 1; i <= maxSeat; i++ {
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
	poolChannel <- struct{}{}
	mu.RLock()
	defer func() {
		mu.RUnlock()
		<-poolChannel
	}()

	room, ok := dataRooms[filmId]
	if !ok {
		return nil, errors.New("Bad film ID, empty room for film")
	}

	return room, nil
}

func InitCurrentUser(id uint, name string) {
	if _, ok := dataUsers[id]; ok {
		return
	}

	initUser(id, name)
}

func AuthUser(name string) string {
	return getAuthToken(name)
}

func IsValidToken(token string) bool {
	poolChannel <- struct{}{}
	mu.RLock()
	defer func() {
		mu.RUnlock()
		<-poolChannel
	}()

	for _, user := range dataUsers {
		if token == getToken(user) {
			return true
		}
	}

	return false
}

func GetFilms() map[uint]string {
	poolChannel <- struct{}{}
	mu.RLock()
	defer func() {
		mu.RUnlock()
		<-poolChannel
	}()

	return films
}

func BuyTicket(filmId, placeId, userId uint) (*Ticket, error) {
	poolChannel <- struct{}{}
	mu.Lock()
	defer func() {
		mu.Unlock()
		<-poolChannel
	}()

	place, ok := dataRooms[filmId].places[placeId]
	if !ok {
		return nil, errors.New("Place ID not found")
	}
	if !place.IsFree() {
		return nil, errors.New("Place ID is is taken")
	}

	place.userId = userId

	ticket := newTicket(filmId, filmId, placeId)
	dataUsers[userId].tickets[ticket.id] = ticket

	return ticket, nil
}

func GetTickets(userId uint) (map[uint]*Ticket, error) {
	poolChannel <- struct{}{}
	mu.RLock()
	defer func() {
		mu.RUnlock()
		<-poolChannel
	}()

	user, ok := dataUsers[userId]
	if !ok {
		return nil, errors.New("User not found")
	}

	return user.tickets, nil
}

func DeleteTicket(userId uint, ticketId uint) error {
	poolChannel <- struct{}{}
	mu.Lock()
	defer func() {
		mu.Unlock()
		<-poolChannel
	}()

	user, ok := dataUsers[userId]
	if !ok {
		return errors.New("User not found")
	}

	ticket, ok := user.tickets[ticketId]
	if !ok {
		return errors.New("Ticket not found")
	}

	dataRooms[ticket.filmId].places[ticket.placeId].userId = 0
	delete(user.tickets, ticketId)

	return nil
}
