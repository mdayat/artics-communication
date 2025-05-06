interface MeetingRoom {
  id: string;
  name: string;
  created_at: string;
}

interface TimeSlot {
  id: string;
  start_date: string;
  end_date: string;
  created_at: string;
}

interface MeetingRoomWithTimeSlotsResponse extends MeetingRoom {
  time_slots: TimeSlot[];
}

export type { MeetingRoomWithTimeSlotsResponse, TimeSlot, MeetingRoom };
