import type { MeetingRoom, TimeSlot } from "./meetingRoom";

interface UserResponse {
  id: string;
  email: string;
  name: string;
  role: string;
  created_at: string;
}

interface RegisterRequest {
  username: string;
  email: string;
  password: string;
}

interface LoginRequest {
  email: string;
  password: string;
}

interface UserReservationResponse {
  id: string;
  meeting_room: MeetingRoom;
  time_slot: TimeSlot;
  canceled: boolean;
  canceled_at: string | null;
  reserved_at: string;
}

export type {
  UserResponse,
  RegisterRequest,
  LoginRequest,
  UserReservationResponse,
};
