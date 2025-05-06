interface CreateReservationRequest {
  meeting_room_id: string;
  time_slot_id: string;
}

interface ReservationResponse {
  id: string;
  user_id: string;
  meeting_room_id: string;
  time_slot_id: string;
  canceled: boolean;
  canceled_at: string | null;
  reserved_at: string;
}

export type { CreateReservationRequest, ReservationResponse };
