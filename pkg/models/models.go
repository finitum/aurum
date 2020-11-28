package models

type Group struct {
	Name              string `json:"name,omitempty"`
	AllowRegistration bool   `json:"allow_registration,omitempty"`
}

type User struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
}

type Role int

const (
	RoleUser Role = iota + 1
	RoleAdmin
)

type GroupWithRole struct {
	Group
	Role Role `json:"role,omitempty"`
}
