package ecc

// Copied from Cloudflare's cfssl library and changed slightly to facilitate our own workflow
// https://github.com/cloudflare/cfssl/blob/6aeb6e3882b8629628efe2e6b4ef89f2ca8ec044/helpers/derhelpers/ed25519.go

// TODO: Rewrite all crypto.*Key and ed25519.*Key to our own ecc.*Key types for cleaner usage and less casting

import (
	"crypto"
	"crypto/ed25519"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
)

var errEd25519WrongID = errors.New("incorrect object identifier")
var errEd25519WrongKeyType = errors.New("incorrect key type")

// ed25519OID is the OID for the Ed25519 signature scheme: see
// https://datatracker.ietf.org/doc/draft-ietf-curdle-pkix-04.
var ed25519OID = asn1.ObjectIdentifier{1, 3, 101, 112}

// subjectPublicKeyInfo reflects the ASN.1 object defined in the X.509 standard.
//
// This is defined in crypto/x509 as "publicKeyInfo".
type subjectPublicKeyInfo struct {
	Algorithm pkix.AlgorithmIdentifier
	PublicKey asn1.BitString
}

// marshalEd25519PublicKey creates a DER-encoded SubjectPublicKeyInfo for an
// ecc public key, as defined in
// https://tools.ietf.org/html/draft-ietf-curdle-pkix-04. This is analagous to
// MarshalPKIXPublicKey in crypto/x509, which doesn't currently support Ed25519.
func marshalEd25519PublicKey(pk crypto.PublicKey) ([]byte, error) {
	pub, ok := pk.(ed25519.PublicKey)
	if !ok {
		return nil, errEd25519WrongKeyType
	}

	spki := subjectPublicKeyInfo{
		Algorithm: pkix.AlgorithmIdentifier{
			Algorithm: ed25519OID,
		},
		PublicKey: asn1.BitString{
			BitLength: len(pub) * 8,
			Bytes:     pub,
		},
	}

	return asn1.Marshal(spki)
}

// ParseEd25519PublicKey returns the Ed25519 public key encoded by the input.
func parseEd25519PublicKey(der []byte) (ed25519.PublicKey, error) {
	var spki subjectPublicKeyInfo
	if rest, err := asn1.Unmarshal(der, &spki); err != nil {
		return nil, err
	} else if len(rest) > 0 {
		return nil, errors.New("SubjectPublicKeyInfo too long")
	}

	if !spki.Algorithm.Algorithm.Equal(ed25519OID) {
		return nil, errEd25519WrongID
	}

	if spki.PublicKey.BitLength != ed25519.PublicKeySize*8 {
		return nil, errors.New("SubjectPublicKeyInfo PublicKey length mismatch")
	}

	return spki.PublicKey.Bytes, nil
}

// oneAsymmetricKey reflects the ASN.1 structure for storing private keys in
// https://tools.ietf.org/html/draft-ietf-curdle-pkix-04, excluding the optional
// fields, which we don't use here.
//
// This is identical to pkcs8 in crypto/x509.
type oneAsymmetricKey struct {
	Version    int
	Algorithm  pkix.AlgorithmIdentifier
	PrivateKey []byte
}

// curvePrivateKey is the inner type of the SecretKey field of
// oneAsymmetricKey.
type curvePrivateKey []byte

// marshalEd25519PrivateKey returns a DER encoding of the input private key as
// specified in https://tools.ietf.org/html/draft-ietf-curdle-pkix-04.
func marshalEd25519PrivateKey(sk crypto.PrivateKey) ([]byte, error) {
	priv, ok := sk.(ed25519.PrivateKey)
	if !ok {
		return nil, errEd25519WrongKeyType
	}

	// Marshal the inner CurvePrivateKey.
	curvePrivateKey, err := asn1.Marshal(priv.Seed())
	if err != nil {
		return nil, err
	}

	// Marshal the OneAsymmetricKey.
	asym := oneAsymmetricKey{
		Version: 0,
		Algorithm: pkix.AlgorithmIdentifier{
			Algorithm: ed25519OID,
		},
		PrivateKey: curvePrivateKey,
	}
	return asn1.Marshal(asym)
}

// parseEd25519PrivateKey returns the Ed25519 private key encoded by the input.
func parseEd25519PrivateKey(der []byte) (ed25519.PrivateKey, error) {
	asym := new(oneAsymmetricKey)
	if rest, err := asn1.Unmarshal(der, asym); err != nil {
		return nil, err
	} else if len(rest) > 0 {
		return nil, errors.New("OneAsymmetricKey too long")
	}

	// Check that the key type is correct.
	if !asym.Algorithm.Algorithm.Equal(ed25519OID) {
		return nil, errEd25519WrongID
	}

	// Unmarshal the inner CurvePrivateKey.
	seed := new(curvePrivateKey)
	if rest, err := asn1.Unmarshal(asym.PrivateKey, seed); err != nil {
		return nil, err
	} else if len(rest) > 0 {
		return nil, errors.New("CurvePrivateKey too long")
	}

	return ed25519.NewKeyFromSeed(*seed), nil
}
