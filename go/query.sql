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

-- name: SelectAvailableMeetingRooms :many
SELECT * FROM meeting_room mr
WHERE EXISTS (
    SELECT 1
    FROM time_slot ts
    WHERE ts.meeting_room_id = mr.id
    AND NOT EXISTS (
        SELECT 1 
        FROM reservation r 
        WHERE r.meeting_room_id = mr.id 
        AND r.time_slot_id = ts.id 
        AND r.status != 'canceled'
    )
);