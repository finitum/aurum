package models

type Application struct {
	Name string `json:"name,omitempty"`
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
