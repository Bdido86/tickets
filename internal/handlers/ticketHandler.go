package handlers

import (
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/storage"
	"gitlab.ozon.dev/Bdido86/movie-tickets/tools"
	"strconv"
)

func ticketsFunc(userId uint) (res string) {
	tickets, err := storage.GetTickets(userId)
	if err != nil {
		return err.Error()
	}

	if len(tickets) == 0 {
		return "У вас нет билетов " + tools.EmojiPensiveFace + " \r\nФильмы на сегодня /films"
	}

	res += "Список ваших билетов:\r\n"
	films := storage.GetFilms()
	for _, ticket := range tickets {
		res += tools.EmojiCheckGreen + " Билет " + strconv.FormatUint(uint64(ticket.GetId()), 10) + ", место " +
			strconv.FormatUint(uint64(ticket.GetPlaceId()), 10) + ", зал " +
			strconv.FormatUint(uint64(ticket.GetRoomId()), 10) + ", фильм '" + films[ticket.GetFilmId()] +
			"'; вернуть " + tools.EmojiRightPoint + " /ticket_" + strconv.FormatUint(uint64(ticket.GetId()), 10) + "\r\n"
	}
	return
}

func ticketFunc(arguments []string, userId uint) (res string) {
	if len(arguments) != 1 {
		return "Bad arguments for ticket return ID"
	}

	ticketId64, err := strconv.ParseUint(arguments[0], 10, 32)
	if err != nil {
		return "Bad arguments for ticket return ID, bad type ID"
	}

	ticketId := uint(ticketId64)
	err = storage.ReturnTicket(userId, ticketId)
	if err != nil {
		return err.Error()
	}

	return tools.EmojiCheckBlack + " Билет номер " + strconv.FormatUint(ticketId64, 10) + " успешно возвращен!"
}
