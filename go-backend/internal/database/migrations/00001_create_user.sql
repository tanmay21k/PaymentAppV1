-- +goose Up
CREATE TABLE IF NOT EXISTS "user" (
    firstname text,
    lastname text,
    username text PRIMARY KEY,
    password text
);

-- +goose Down
DROP TABLE IF EXISTS "user";