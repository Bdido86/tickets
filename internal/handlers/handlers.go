package handlers

import "ozon/go-hw-bot/internal/commander"

const (
	startCmd = "start"
	helpCmd  = "help"
	filmsCmd = "films"
)

func Run(m commander.Message) string {
	switch m.Cmd() {
	case startCmd:
		return startFunc(m.UserName())
	case helpCmd:
		return helpFunc()
		//case FilmsCmd:
		//	return filmsFunc()
	}

	return "Неизвестная команда. Посмотрите справку по командам /help"
}

func AddHandler(c *commander.Commander) {
	c.RegisterHandler(Run)
}
