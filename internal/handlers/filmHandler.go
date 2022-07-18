package handlers

import (
	"ozon/go-hw-bot/internal/storage"
	"strconv"
	"time"
)

func filmsFunc() (res string) {
	films := storage.GetFilms()
	dt := time.Now()

	res = "Расписание на сегодня " + dt.Format("02-01-2006") + " \xF0\x9F\x8E\xA5: \r\n\r\n"
	for id, film := range films {
		res += "\xE2\x9C\x94 " + film + " \xF0\x9F\x91\x89 /film_" + strconv.Itoa(id) + "\r\n"
	}
	return
}

func filmFunc(arguments []string) (res string) {
	if len(arguments) != 1 {
		return "Bad film select ID"
	}

	filmId, err := strconv.Atoi(arguments[0])
	if err != nil {
		return "Bad film select ID"
	}

	films := storage.GetFilms()
	film, ok := films[filmId]
	if !ok {
		return "Bad film select ID"
	}

	res = "Ваш выбор фильма\r\n\xE2\x9C\x85 " + film + "!\r\nПоздравляем! А теперь надо выбрать место в зале:\r\n"

	return
}
