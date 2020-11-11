package models

type AccessResponse struct {
	ApplicationName string
	Username        string
	AllowedAccess   bool
	Role            Role
}

type PublicKeyResponse struct {
	PublicKey string `json:"public_key"`
}
