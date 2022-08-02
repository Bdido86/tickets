-- +goose Up
-- +goose StatementBegin
INSERT INTO public.films (name)
values ('Один дома1'),
       ('Один дома2'),
       ('Один дома3'),
       ('Зеленая миля'),
       ('Бойцовский клуб'),
       ('Джентельмены удачи'),
       ('Титаник'),
       ('Вверх (мультфильм)'),
       ('Один дома 2'),
       ('Курсы GO в ozon{ech (обучение)'),
       ('Крёстный отец'),
       ('Тёмный рыцарь'),
       ('Начало'),
       ('Славные парни'),
       ('Семь'),
       ('Пианист');


INSERT INTO public.rooms (count_places)
values (50),
       (40),
       (40),
       (30),
       (25),
       (20),
       (10),
       (5);


INSERT INTO public.film_room (film_id, room_id)
values (1, 1),
       (2, 2),
       (3, 3),
       (4, 4),
       (5, 5),
       (6, 6),
       (7, 7),
       (8, 8),
       (9, 1),
       (10, 2),
       (11, 3),
       (12, 4),
       (13, 5),
       (14, 6),
       (15, 7),
       (16, 8);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE public.films;
TRUNCATE TABLE public.rooms;
TRUNCATE TABLE public.film_room;
-- +goose StatementEnd
