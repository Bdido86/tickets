package handlers

func helpFunc() string {
	return `
Доступные команды для выбора фильма
/start - начать диалог заново
/help - полный список команд, вы тут :)
----
/films - список фильмов на сегодня
/film_<filmId> - выбор фильма для просмотра
/film_<filmId>_<placeId> - выбор места в зале для просмотра
----
`
}
