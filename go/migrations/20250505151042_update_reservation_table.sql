-- Drop index "idx_active_reservation" from table: "reservation"
DROP INDEX "idx_active_reservation";
-- Modify "reservation" table
ALTER TABLE "reservation" DROP CONSTRAINT "reservation_status_check", DROP CONSTRAINT "reservation_meeting_room_id_time_slot_id_key", DROP COLUMN "status", ADD COLUMN "canceled" boolean NOT NULL DEFAULT false, ADD COLUMN "canceled_at" timestamptz NULL;
-- Create index "idx_active_reservation" to table: "reservation"
CREATE UNIQUE INDEX "idx_active_reservation" ON "reservation" ("meeting_room_id", "time_slot_id") WHERE (NOT canceled);
