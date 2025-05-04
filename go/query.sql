-- name: SelectUser :one
SELECT * FROM "user" WHERE id = $1;