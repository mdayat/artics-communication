package dtos

type CanceledReservationResponse struct {
	Id            string  `json:"id"`
	UserId        string  `json:"user_id"`
	MeetingRoomId string  `json:"meeting_room_id"`
	TimeSlotId    string  `json:"time_slot_id"`
	Canceled      bool    `json:"canceled"`
	CanceledAt    *string `json:"canceled_at"`
	ReservedAt    string  `json:"reserved_at"`
}
