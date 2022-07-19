package handlers

import "gitlab.ozon.dev/Bdido86/movie-tickets/internal/commander"

const (
	startCmd = "start"
	helpCmd  = "help"

	filmCmd  = "film"
	filmsCmd = "films"

	ticketCmd  = "ticket"
	ticketsCmd = "tickets"
)

func Run(m commander.Message) string {
	initCurrentUser(m.UserId(), m.UserName())

	switch m.Cmd() {
	case startCmd:
		return startFunc(m.UserName())
	case helpCmd:
		return helpFunc()

	case filmCmd:
		return filmFunc(m.Arguments(), m.UserId())
	case filmsCmd:
		return filmsFunc()

	case ticketCmd:
		return ticketFunc(m.Arguments(), m.UserId())
	case ticketsCmd:
		return ticketsFunc(m.UserId())
	}

	return "Неизвестная команда. Справка по командам /help"
}

func AddHandler(c *commander.Commander) {
	c.RegisterHandler(Run)
}
