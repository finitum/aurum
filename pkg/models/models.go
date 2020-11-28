package models

type User struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
	Role     Role   `json:"role"`
}

type Role int

const (
	RoleUser Role = iota + 1
	RoleAdmin
)
