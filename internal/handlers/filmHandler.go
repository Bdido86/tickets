package handlers

import (
	"gitlab.ozon.dev/Bdido86/go-hw-bot/internal/storage"
	"strconv"
	"time"
)

var seatsEmoji = map[uint]string{
	0: "\x39\xE2\x83\xA3",
	1: "\x31\xE2\x83\xA3",
	2: "\x32\xE2\x83\xA3",
	3: "\x33\xE2\x83\xA3",
	4: "\x34\xE2\x83\xA3",
	5: "\x35\xE2\x83\xA3",
	6: "\x36\xE2\x83\xA3",
	7: "\x37\xE2\x83\xA3",
	8: "\x38\xE2\x83\xA3",
	9: "\x39\xE2\x83\xA3",
}

func filmsFunc() (res string) {
	films := storage.GetFilms()
	dt := time.Now()

	res = "Расписание на сегодня " + dt.Format("02-01-2006") + " \xF0\x9F\x8E\xA5: \r\n\r\n"
	for id, film := range films {
		res += "\xE2\x9C\x94 " + film + " \xF0\x9F\x91\x89 /film_" + strconv.FormatUint(uint64(id), 10) + "\r\n"
	}
	return
}

func filmFunc(arguments []string, userId uint) (res string) {
	if len(arguments) == 0 || len(arguments) > 2 {
		return "Bad arguments for film"
	}

	filmId64, err := strconv.ParseUint(arguments[0], 10, 32)
	if err != nil {
		return "Bad film select ID, bad type ID"
	}

	filmId := uint(filmId64)
	films := storage.GetFilms()
	film, ok := films[filmId]
	if !ok {
		return "Bad film select ID, not found film"
	}

	room, err := storage.GetRoom(filmId)
	if err != nil {
		return err.Error()
	}

	places := room.GetPlaces()
	res = "Ваш выбор фильма\r\n\xE2\x9C\x85 " + film + ":\r\n"
	res += "\xE2\x9D\x97 Показ будет идти в зале номер " + strconv.FormatUint(uint64(room.GetNumber()), 10) + "\r\n"

	if len(arguments) == 1 {
		res += "Необходимо выбрать место:\r\n"
		for _, place := range places {
			res += seatsEmoji[place.GetNumber()] + ": "
			if place.IsFree() {
				res += "выбрать \xF0\x9F\x91\x89 /film_" + strconv.FormatUint(filmId64, 10) + "_" + strconv.FormatUint(uint64(place.GetNumber()), 10)
			} else {
				if place.GetUserId() == userId {
					res += "вы ранее купили это место"
				} else {
					res += "уже продано"
				}
			}
			res += "\r\n"
		}

		return
	}

	placeId64, err := strconv.ParseUint(arguments[1], 10, 32)
	if err != nil {
		return "Bad place select ID, bad type ID"
	}

	placeId := uint(placeId64)
	place, ok := places[placeId]
	if !ok {
		return "Bad place select ID, place not found"
	}

	if !place.IsFree() {
		if place.GetUserId() == userId {
			return "Bad place select ID, вы ранее уже купили это место"
		} else {
			return "Bad place select ID, место уже занято"
		}
	}

	ticket, err := storage.BuyTicket(filmId, placeId, userId)
	if err != nil {
		return err.Error()
	}

	res += "Место номер " + strconv.FormatUint(uint64(ticket.GetPlaceId()), 10) + "\r\n"
	res += "Номер билета " + strconv.FormatUint(uint64(ticket.GetId()), 10) + "\r\n\r\n"
	res += "ПРИЯТНОГО ПРОСМОТРА"
	return
}
