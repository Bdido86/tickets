package postgres

const (
	SelectAllTablesWithoutGoose = `SELECT table_name FROM information_schema.tables 
									WHERE table_schema='public' AND table_type='BASE TABLE' AND table_name not like '%goose_%';`

	InsertTicket = `INSERT INTO tickets(user_id, film_id, room_id, place)
				    VALUES($1,$2,$3,$4);`

	InsertUser = `INSERT INTO users(name, token)
				    VALUES($1,$2);`

	InsertFilm = `INSERT INTO films(name)
				    VALUES($1);`

	InsertRoom = `INSERT INTO rooms(count_places)
				    VALUES($1);`

	InsertFilmRoom = `INSERT INTO film_room(film_id, room_id)
				    VALUES($1,$2);`
)
