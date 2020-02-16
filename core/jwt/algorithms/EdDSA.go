// By Blain Smith
// From https://blainsmith.com/articles/signing-jwts-with-gos-crypto-ed25519/

package algorithms

import (
	"encoding/asn1"
	"encoding/pem"
	"errors"
)

import (
	"crypto/ed25519"
	"github.com/dgrijalva/jwt-go"
)

var (
	// Sadly this is missing from crypto/ecdsa compared to crypto/rsa
	ErrEdDSAVerification = errors.New("ed25519: verification error")
)

func init() {
	var edDSASigningMethod SigningMethodEdDSA
	jwt.RegisterSigningMethod(edDSASigningMethod.Alg(), func() jwt.SigningMethod { return &edDSASigningMethod })
}

type SigningMethodEdDSA struct{}

func (m *SigningMethodEdDSA) Alg() string {
	return "EdDSA"
}

func (m *SigningMethodEdDSA) Verify(signingString string, signature string, key interface{}) error {
	var err error

	var sig []byte
	if sig, err = jwt.DecodeSegment(signature); err != nil {
		return err
	}

	var ed25519Key ed25519.PublicKey
	var ok bool
	if ed25519Key, ok = key.(ed25519.PublicKey); !ok {
		return jwt.ErrInvalidKeyType
	}

	if len(ed25519Key) != ed25519.PublicKeySize {
		return jwt.ErrInvalidKey
	}

	if ok := ed25519.Verify(ed25519Key, []byte(signingString), sig); !ok {
		return ErrEdDSAVerification
	}

	return nil
}

func (m *SigningMethodEdDSA) Sign(signingString string, key interface{}) (str string, err error) {
	var ed25519Key ed25519.PrivateKey
	var ok bool
	if ed25519Key, ok = key.(ed25519.PrivateKey); !ok {
		return "", jwt.ErrInvalidKeyType
	}

	if len(ed25519Key) != ed25519.PrivateKeySize {
		return "", jwt.ErrInvalidKey
	}

	// Sign the string and return the encoded result
	sig := ed25519.Sign(ed25519Key, []byte(signingString))
	return jwt.EncodeSegment(sig), nil
}


type ed25519PrivKey struct {
	Version          int
	ObjectIdentifier struct {
		ObjectIdentifier asn1.ObjectIdentifier
	}
	PrivateKey []byte
}

func DecodeEdDSAPrivatekey (privateKeyPEM []byte) (ed25519.PrivateKey, error) {

	var block *pem.Block
	block, _ = pem.Decode(privateKeyPEM)

	var asn1PrivKey ed25519PrivKey
	_, err := asn1.Unmarshal(block.Bytes, &asn1PrivKey)
	if err != nil {
		return nil, err
	}

	privateKey := ed25519.NewKeyFromSeed(asn1PrivKey.PrivateKey[2:])

	return privateKey, nil
}

type ed25519PubKey struct {
	OBjectIdentifier struct {
		ObjectIdentifier asn1.ObjectIdentifier
	}
	PublicKey asn1.BitString
}

func DecodeEdDSAPublicKey (publicKeyPEM []byte) (ed25519.PublicKey, error) {
	var block *pem.Block
	block, _ = pem.Decode(publicKeyPEM)

	var asn1PubKey ed25519PubKey
	_, err := asn1.Unmarshal(block.Bytes, &asn1PubKey)

	if err != nil {
		return nil, err
	}

	publicKey := ed25519.PublicKey(asn1PubKey.PublicKey.Bytes)

	return publicKey, nil
}