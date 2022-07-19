package handlers

import (
	"gitlab.ozon.dev/Bdido86/movie-tickets/internal/storage"
	"gitlab.ozon.dev/Bdido86/movie-tickets/tools"
)

func startFunc(userName string) string {
	return "Привет " + userName + " " + tools.EmojiWinkingFace + "\n\r" +
		"Сегодня хороший день, чтобы сходить в кино.\r\n\r\n" +
		"Для выбора фильма жми /films \n\r" +
		"Для справочной информации используй команду /help"
}

func initCurrentUser(userId uint, userName string) {
	storage.InitCurrentUser(userId, userName)
}
