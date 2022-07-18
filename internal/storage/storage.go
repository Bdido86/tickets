package storage

var films = map[int]string{
	1: "Зеленая миля",
	2: "Бойцовский клуб",
	3: "Джентльмены удачи",
	4: "Титаник",
	5: "Вверх (мультфильм)",
}

func GetFilms() map[int]string {
	return films
}
