package models

import "github.com/google/uuid"

type Application struct {
	AppId uuid.UUID `json:"appID,omitempty"`
	Name  string    `json:"name,omitempty"`
}

type User struct {
	Username string    `json:"username,omitempty"`
	Password string    `json:"password,omitempty"`
	Email    string    `json:"email,omitempty"`
}

type Role int

const (
	RoleUser Role = iota+1
	RoleAdmin
)


