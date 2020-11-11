package models

import "github.com/google/uuid"

type Application struct {
	AppId uuid.UUID `json:"appID"`
	Name  string    `json:"name"`
}

type User struct {
	UserId   uuid.UUID `json:"userID"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Email    string    `json:"email"`
}

type Role int

const (
	RoleUser Role = iota
	RoleAdmin
)
