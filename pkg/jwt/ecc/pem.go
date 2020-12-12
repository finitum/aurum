package ecc

import (
	"crypto/ed25519"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

const (
	publicKeyPemHeader        = "PUBLIC KEY"
	ed25519SecretKeyPemHeader = "ED25519 PRIVATE KEY"
)

type Key interface {
	// Converts a key to a Pem string
	ToPem() (string, error)

	// Writes a pem string to a given filepath
	WriteToFile(path string) error
}

func (k SecretKey) GetPublicKey() PublicKey {
	return PublicKey(ed25519.PrivateKey(k).Public().(ed25519.PublicKey))
}

func (k PublicKey) ToPem() (string, error) {
	return toPem(k, true)
}

func (k SecretKey) ToPem() (string, error) {
	return toPem(k, false)
}

func (k PublicKey) WriteToFile(path string) error {
	return writeToFile(k, path)
}

func (k SecretKey) WriteToFile(path string) error {
	return writeToFile(k, path)
}

// Matches finds whether or not a public key belongs to this private key
func (k SecretKey) Matches(key PublicKey) bool {
	return ed25519.PublicKey(k.GetPublicKey()).Equal(ed25519.PublicKey(key))
}

func writeToFile(k Key, path string) error {
	kpem, err := k.ToPem()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, []byte(kpem), 0o600)
}

func toPem(key []byte, public bool) (string, error) {
	var Type string

	var pkey []byte
	var err error

	if public {
		Type = publicKeyPemHeader
		pkey, err = marshalEd25519PublicKey(ed25519.PublicKey(key))
	} else {
		Type = ed25519SecretKeyPemHeader
		pkey, err = marshalEd25519PrivateKey(ed25519.PrivateKey(key))
	}

	if err != nil {
		return "", err
	}

	block := pem.Block{
		Type:  Type,
		Bytes: pkey[:],
	}

	bytes := pem.EncodeToMemory(&block)

	return string(bytes[:]), nil
}

// returns either a secret or public key based on the pem
func FromPem(data []byte) (k Key, err error) {
	dec, _ := pem.Decode(data)
	if dec == nil {
		return nil, errors.New("couldn't decode pem key")
	}

	switch dec.Type {
	case publicKeyPemHeader:
		// public key
		der := dec.Bytes
		pk, err := parseEd25519PublicKey(der)
		if err != nil {
			return nil, err
		}
		k = PublicKey(pk)

	case ed25519SecretKeyPemHeader:
		// secret key
		der := dec.Bytes
		sk, err := parseEd25519PrivateKey(der)
		if err != nil {
			return nil, err
		}
		k = SecretKey(sk)

	default:
		// unknown key
		err = errors.New("unknown key type")
	}

	return
}

func FromFile(path string) (Key, error) {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	println(string(key))

	return FromPem(key)
}
