package handlers

import (
	"ozon/go-hw-bot/internal/storage"
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
	if len(arguments) != 1 {
		return "Bad film select ID"
	}

	filmId, err := strconv.ParseUint(arguments[0], 10, 32)
	if err != nil {
		return "Bad film select ID"
	}

	films := storage.GetFilms()
	film, ok := films[uint(filmId)]
	if !ok {
		return "Bad film select ID"
	}

	room, err := storage.GetRoom(uint(filmId))
	if err != nil {
		return err.Error()
	}

	res = "Ваш выбор фильма\r\n\xE2\x9C\x85 " + film + ":\r\n"
	res += "\xE2\x9D\x97 Показ будет идти в зале номер " + strconv.FormatUint(uint64(room.GetNumber()), 10) + "\r\n"
	res += "Необходимо выбрать место \r\n"

	places := room.GetPlaces()
	for _, place := range places {
		res += seatsEmoji[place.GetNumber()] + ": "
		if place.IsFree() {
			res += "выбрать \xF0\x9F\x91\x89 /film_" + strconv.FormatUint(uint64(filmId), 10) + "_" + strconv.FormatUint(uint64(place.GetNumber()), 10)
		} else {
			if place.GetUserId() == userId {
				res += "вы ранее купили это место"
				return
			} else {
				res += "уже продано"
			}
		}
		res += "\r\n"
	}

	return
}
