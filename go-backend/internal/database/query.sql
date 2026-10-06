-- name: CreateUser :one
INSERT INTO "user" (firstname, lastname, username, password)
VALUES ($1, $2, $3, $4)
RETURNING *;
