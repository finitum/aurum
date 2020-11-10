// By Blain Smith
// From https://blainsmith.com/articles/signing-jwts-with-gos-crypto-ed25519/
// Modified in places to use our own keytype

package ecc

import (
	"crypto/ed25519"
	"encoding/asn1"
	"errors"
	"github.com/dgrijalva/jwt-go"
)

var (
	// Sadly this is missing from crypto/ecdsa compared to crypto/rsa
	ErrEdDSAVerification = errors.New("ecc: verification error")
)

func init() {
	var edDSASigningMethod SigningMethodEdDSA
	jwt.RegisterSigningMethod(edDSASigningMethod.Alg(), func() jwt.SigningMethod { return &edDSASigningMethod })
}

type SigningMethodEdDSA struct{}

type OBjectIdentifier struct {
	ObjectIdentifier asn1.ObjectIdentifier
}

type Ed25519PrivKey struct {
	Version          int
	OBjectIdentifier OBjectIdentifier
	PrivateKey       []byte
}

type Ed25519PubKey struct {
	OBjectIdentifier OBjectIdentifier
	PublicKey        asn1.BitString
}

func (m *SigningMethodEdDSA) Alg() string {
	return "EdDSA"
}

func (m *SigningMethodEdDSA) Verify(signingString string, signature string, key interface{}) error {
	var err error

	var sig []byte
	if sig, err = jwt.DecodeSegment(signature); err != nil {
		return err
	}

	var ed25519Key PublicKey
	var ok bool
	if ed25519Key, ok = key.(PublicKey); !ok {
		return jwt.ErrInvalidKeyType
	}

	if len(ed25519Key) != ed25519.PublicKeySize {
		return jwt.ErrInvalidKey
	}

	if ok := ed25519.Verify(ed25519.PublicKey(ed25519Key), []byte(signingString), sig); !ok {
		return ErrEdDSAVerification
	}

	return nil
}

func (m *SigningMethodEdDSA) Sign(signingString string, key interface{}) (str string, err error) {
	var ed25519Key SecretKey
	var ok bool
	if ed25519Key, ok = key.(SecretKey); !ok {
		return "", jwt.ErrInvalidKeyType
	}

	if len(ed25519Key) != ed25519.PrivateKeySize {
		return "", jwt.ErrInvalidKey
	}

	// Sign the string and return the encoded result
	sig := ed25519.Sign(ed25519.PrivateKey(ed25519Key), []byte(signingString))
	return jwt.EncodeSegment(sig), nil
}
