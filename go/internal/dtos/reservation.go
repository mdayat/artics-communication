package dtos

type CreateReservationRequest struct {
	MeetingRoomId string `json:"meeting_room_id"`
	TimeSlotId    string `json:"time_slot_id"`
}

type ReservationResponse struct {
	Id            string  `json:"id"`
	UserId        string  `json:"user_id"`
	MeetingRoomId string  `json:"meeting_room_id"`
	TimeSlotId    string  `json:"time_slot_id"`
	Canceled      bool    `json:"canceled"`
	CanceledAt    *string `json:"canceled_at"`
	ReservedAt    string  `json:"reserved_at"`
}

type EnrichedReservationResponse struct {
	Id          string       `json:"id"`
	User        UserResponse `json:"user"`
	MeetingRoom MeetingRoom  `json:"meeting_room"`
	TimeSlot    TimeSlot     `json:"time_slot"`
	Canceled    bool         `json:"canceled"`
	CanceledAt  *string      `json:"canceled_at"`
	ReservedAt  string       `json:"reserved_at"`
}
