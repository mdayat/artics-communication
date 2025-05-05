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

type UserReservationResponse struct {
	Id          string      `json:"id"`
	MeetingRoom MeetingRoom `json:"meeting_room"`
	TimeSlot    TimeSlot    `json:"time_slot"`
	Canceled    bool        `json:"canceled"`
	CanceledAt  *string     `json:"canceled_at"`
	ReservedAt  string      `json:"reserved_at"`
}
