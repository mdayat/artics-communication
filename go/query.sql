-- name: InsertUser :one
INSERT INTO "user" (id, email, password, name, role)
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: SelectUser :one
SELECT * FROM "user" WHERE id = $1;

-- name: SelectUserByEmail :one
SELECT * FROM "user" WHERE email = $1;

-- name: SelectUserReservations :many
SELECT r.*, sqlc.embed(mr), sqlc.embed(ts) FROM reservation r
  JOIN meeting_room mr ON mr.id = r.meeting_room_id
  JOIN time_slot ts ON ts.id = r.time_slot_id
WHERE user_id = $1;