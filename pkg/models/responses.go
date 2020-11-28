package models

type AccessStatus struct {
	GroupName string
	Username        string
	AllowedAccess   bool
	Role            Role
}

type PublicKeyResponse struct {
	PublicKey string `json:"public_key"`
}
