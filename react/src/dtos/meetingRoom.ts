interface TimeSlot {
  id: string;
  start_date: string;
  end_date: string;
  created_at: string;
}

interface MeetingRoomWithTimeSlotsResponse {
  id: string;
  name: string;
  created_at: string;
  time_slots: TimeSlot[];
}

export type { MeetingRoomWithTimeSlotsResponse, TimeSlot };
