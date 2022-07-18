package handlers

import "ozon/go-hw-bot/internal/commander"

const (
	startCmd = "start"
	helpCmd  = "help"

	filmsCmd = "films"
	filmCmd  = "film"
)

func Run(m commander.Message) string {
	switch m.Cmd() {
	case startCmd:
		return startFunc(m.UserId(), m.UserName())
	case helpCmd:
		return helpFunc()

	case filmCmd:
		return filmFunc(m.Arguments(), m.UserId())
	case filmsCmd:
		return filmsFunc()
	}

	return "Неизвестная команда. Посмотрите справку по командам /help"
}

func AddHandler(c *commander.Commander) {
	c.RegisterHandler(Run)
}
