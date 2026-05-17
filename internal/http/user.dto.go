package http

type UserBasicProfile struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}
