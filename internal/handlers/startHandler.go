package handlers

import "gitlab.ozon.dev/Bdido86/go-hw-bot/internal/storage"

func startFunc(userName string) string {
	return "Привет " + userName + "! \n\r" +
		"Сегодня хороший день, чтобы сходить в кино. Для выбора фильма жми /films.\n\r" +
		"Для справочной информации используй команду /help"
}

func initCurrentUser(userId uint, userName string) {
	storage.InitCurrentUser(userId, userName)
}
