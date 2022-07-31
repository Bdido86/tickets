-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.films (
    id  integer NOT NULL PRIMARY KEY,
    name varchar(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS public.users (
    id  integer NOT NULL PRIMARY KEY,
    name varchar(255) NOT NULL UNIQUE,
    token varchar(255) NOT NULL UNIQUE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.films;
DROP TABLE IF EXISTS public.users;
-- +goose StatementEnd
