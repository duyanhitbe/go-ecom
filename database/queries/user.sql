-- name: CreateUser :one
INSERT INTO users (username, password)
VALUES (sqlc.arg('username'), sqlc.arg('password'))
RETURNING *;

-- name: FindUser :many
SELECT *
FROM users
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');

-- name: CountUser :one
SELECT COUNT(*)::int4
FROM users;

-- name: FindOneUserByUsername :one
SELECT *
FROM users
WHERE username = $1
LIMIT 1;