package dtos

const (
	AdminRole string = "admin"
	UserRole  string = "user"
)

type UserResponse struct {
	Id        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

type UserReservation struct {
	Id          string      `json:"id"`
	MeetingRoom MeetingRoom `json:"meeting_room"`
	TimeSlot    TimeSlot    `json:"time_slot"`
	Status      string      `json:"status"`
	ReservedAt  string      `json:"reserved_at"`
}
