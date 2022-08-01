-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.users
(
    id    integer      PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    name  varchar(255) NOT NULL UNIQUE,
    token varchar(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS public.films
(
    id   integer      PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    name varchar(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS public.rooms
(
    id          integer PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    count_seats integer NOT NULL
);
comment on column rooms.count_seats is 'Кол-во мест в зале';

CREATE TABLE IF NOT EXISTS public.film_room
(
    film_id integer NOT NULL,
    room_id integer NOT NULL
);
CREATE UNIQUE INDEX film_room_unique ON public.film_room(film_id, room_id);

CREATE TABLE IF NOT EXISTS public.tickets
(
    id      integer PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    user_id integer NOT NULL,
    film_id integer NOT NULL,
    room_id integer NOT NULL,
    seat    integer NOT NULL
);
CREATE UNIQUE INDEX tickets_user_film_room_seat ON public.tickets(user_id, film_id, room_id, seat);
CREATE INDEX tickets_user ON public.tickets(user_id);
CREATE INDEX tickets_film_room ON public.tickets(film_id, room_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.tickets;
DROP TABLE IF EXISTS public.film_room;
DROP TABLE IF EXISTS public.rooms;
DROP TABLE IF EXISTS public.films;
DROP TABLE IF EXISTS public.users;
-- +goose StatementEnd
