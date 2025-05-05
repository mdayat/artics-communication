package dtos

type MeetingRoom struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

type TimeSlot struct {
	Id        string `json:"id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	CreatedAt string `json:"created_at"`
}

type MeetingRoomWithTimeSlots struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	TimeSlots []TimeSlot
}
