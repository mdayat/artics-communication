-- Create index "idx_active_reservation" to table: "reservation"
CREATE UNIQUE INDEX "idx_active_reservation" ON "reservation" ("meeting_room_id", "time_slot_id") WHERE ((status)::text <> 'canceled'::text);
