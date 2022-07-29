package command

import (
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/pkg/storage"
	"gitlab.ozon.dev/Bdido86/movie-tickets/tools"
	"sort"
	"strconv"
	"time"
)

func filmsFunc() (res string) {
	films := storage.GetFilms()
	dt := time.Now()

	res = "Расписание на сегодня " + dt.Format("02-01-2006") + " " + tools.EmojiCamera + ": \r\n\r\n"
	for id, film := range films {
		res += tools.EmojiCheckBlack + " " + film + " " + tools.EmojiRightPoint + " /film_" + strconv.FormatUint(uint64(id), 10) + "\r\n"
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
	res = "Ваш выбор фильма\r\n" + tools.EmojiCheckGreen + " " + film + ":\r\n"
	res += "Показ будет идти в зале номер " + strconv.FormatUint(uint64(room.GetNumber()), 10) + "\r\n"

	if len(arguments) == 1 {
		res += "\r\nНеобходимо выбрать место:\r\n"

		var keys []int
		for _, place := range places {
			keys = append(keys, int(place.GetNumber()))
		}
		sort.Ints(keys)

		for _, k := range keys {
			place := places[uint(k)]
			res += tools.EmojiNumbers[place.GetNumber()] + ": "
			if place.IsFree() {
				res += "выбрать " + tools.EmojiRightPoint + " /film_" + strconv.FormatUint(filmId64, 10) + "_" + strconv.FormatUint(uint64(place.GetNumber()), 10)
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
	res += "Выбрать фильм " + tools.EmojiRightPoint + " /films \r\n"
	res += "Список билетов " + tools.EmojiRightPoint + " /tickets \r\n"
	return
}
