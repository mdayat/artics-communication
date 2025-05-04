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
