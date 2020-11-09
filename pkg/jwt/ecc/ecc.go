package ecc

import (
	"crypto/ed25519"
	"crypto/rand"
)

// Generates a pair of ed25519 keys and wraps them into the ecc types
func GenerateKey() (PublicKey, SecretKey, error) {
	pk, sk, err := ed25519.GenerateKey(rand.Reader)
	return PublicKey(pk), SecretKey(sk), err
}
