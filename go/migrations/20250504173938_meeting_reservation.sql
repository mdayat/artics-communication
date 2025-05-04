-- Create "meeting_room" table
CREATE TABLE "meeting_room" (
  "id" uuid NOT NULL,
  "name" character varying(255) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);
-- Create "time_slot" table
CREATE TABLE "time_slot" (
  "id" uuid NOT NULL,
  "meeting_room_id" uuid NOT NULL,
  "start_date" timestamptz NOT NULL,
  "end_date" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_time_slot_meeting_room_id" FOREIGN KEY ("meeting_room_id") REFERENCES "meeting_room" ("id") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "time_slot_dates_check" CHECK (start_date < end_date)
);
-- Create "reservation" table
CREATE TABLE "reservation" (
  "id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "meeting_room_id" uuid NOT NULL,
  "time_slot_id" uuid NOT NULL,
  "status" character varying(50) NOT NULL DEFAULT 'in_progress',
  "reserved_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "reservation_meeting_room_id_time_slot_id_key" UNIQUE ("meeting_room_id", "time_slot_id"),
  CONSTRAINT "fk_reservation_meeting_room_id" FOREIGN KEY ("meeting_room_id") REFERENCES "meeting_room" ("id") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "fk_reservation_time_slot_id" FOREIGN KEY ("time_slot_id") REFERENCES "time_slot" ("id") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "fk_reservation_user_id" FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "reservation_status_check" CHECK ((status)::text = ANY ((ARRAY['confirmed'::character varying, 'canceled'::character varying, 'in_progress'::character varying, 'completed'::character varying])::text[]))
);
