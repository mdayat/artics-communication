-- name: BulkInsertUser :copyfrom
INSERT INTO "user" (id, email, password, name, role)
VALUES ($1, $2, $3, $4, $5);

-- name: InsertUser :one
INSERT INTO "user" (id, email, password, name, role)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: SelectUser :one
SELECT * FROM "user" WHERE id = $1;

-- name: SelectUserByEmail :one
SELECT * FROM "user" WHERE email = $1;

-- name: SelectUserReservations :many
SELECT
  r.*,
  sqlc.embed(mr),
  sqlc.embed(ts)
FROM
  reservation r
JOIN
  meeting_room mr ON mr.id = r.meeting_room_id
JOIN
  time_slot ts ON ts.id = r.time_slot_id
WHERE
  user_id = $1;

-- name: CancelUserReservation :one
UPDATE reservation
SET
  canceled = TRUE
WHERE
  id = $1
  AND user_id = $2
RETURNING *;

-- name: BulkInsertMeetingRoom :copyfrom
INSERT INTO meeting_room (id, name)
VALUES ($1, $2);

-- name: SelectMeetingRooms :many
SELECT
  mr.*,
  sqlc.embed(ts)
FROM
  meeting_room mr
JOIN
  time_slot ts ON ts.meeting_room_id = mr.id
ORDER BY
  mr.name;

-- name: SelectAvailableMeetingRooms :many
SELECT
  mr.*,
  sqlc.embed(ts)
FROM
  meeting_room mr
JOIN
  time_slot ts ON ts.meeting_room_id = mr.id
WHERE
  NOT EXISTS (
    SELECT 1 FROM reservation r 
    WHERE
      r.meeting_room_id = mr.id 
      AND r.time_slot_id = ts.id 
      AND r.canceled == FALSE
  )
ORDER BY
  mr.name;

-- name: BulkInsertReservation :copyfrom
INSERT INTO reservation (id, user_id, meeting_room_id, time_slot_id, canceled, canceled_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: SelectReservations :many
SELECT
  r.*,
  sqlc.embed(u),
  sqlc.embed(mr),
  sqlc.embed(ts)
FROM
  reservation r
JOIN
  "user" u ON u.id = r.user_id
JOIN
  meeting_room mr ON mr.id = r.meeting_room_id
JOIN
  time_slot ts ON ts.id = r.time_slot_id;

-- name: CancelReservation :one
UPDATE reservation
SET
  canceled = TRUE
WHERE
  id = $1
RETURNING *;

-- name: BulkInsertTimeSlot :copyfrom
INSERT INTO time_slot (id, meeting_room_id, start_date, end_date)
VALUES ($1, $2, $3, $4);