CREATE TABLE "user" (
  id UUID PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  name VARCHAR(255) NOT NULL,
  role VARCHAR(255) NOT NULL CHECK (role IN ('admin', 'user')),
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE meeting_room (
  id UUID PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE time_slot (
  id UUID PRIMARY KEY,
  meeting_room_id UUID NOT NULL,
  start_date TIMESTAMPTZ NOT NULL,
  end_date TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
  
  CONSTRAINT time_slot_dates_check CHECK (start_date < end_date),

  CONSTRAINT fk_time_slot_meeting_room_id
    FOREIGN KEY (meeting_room_id)
    REFERENCES meeting_room(id)
    ON UPDATE CASCADE
    ON DELETE CASCADE
);

CREATE TABLE reservation (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL,
  meeting_room_id UUID NOT NULL,
  time_slot_id UUID NOT NULL,
  canceled BOOLEAN DEFAULT FALSE NOT NULL,
  canceled_at TIMESTAMPTZ NULL,
  reserved_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

  CONSTRAINT fk_reservation_user_id
    FOREIGN KEY (user_id)
    REFERENCES "user"(id)
    ON UPDATE CASCADE
    ON DELETE CASCADE,

  CONSTRAINT fk_reservation_meeting_room_id
    FOREIGN KEY (meeting_room_id)
    REFERENCES meeting_room(id)
    ON UPDATE CASCADE
    ON DELETE CASCADE,

  CONSTRAINT fk_reservation_time_slot_id
    FOREIGN KEY (time_slot_id)
    REFERENCES time_slot(id)
    ON UPDATE CASCADE
    ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_active_reservation
ON reservation(meeting_room_id, time_slot_id) 
WHERE NOT canceled;