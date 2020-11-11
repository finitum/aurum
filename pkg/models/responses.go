package models

import "github.com/google/uuid"

type AccessResponse struct {
	ApplicationID uuid.UUID
	Username      string
	AllowedAccess bool
	Role          Role
}

type PublicKeyResponse struct {
	PublicKey string `json:"public_key"`
}
