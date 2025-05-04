-- name: InsertUser :one
INSERT INTO "user" (id, email, password, name, role)
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: SelectUser :one
SELECT * FROM "user" WHERE id = $1;

-- name: SelectUserByEmail :one
SELECT * FROM "user" WHERE email = $1;