package ecc

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
)

const (
	generateKeysSTDOUT = "stdout"
	generateKeysFile   = "file"
	generateKeysBoth   = "both"
)

// Generates a pair of ed25519 keys and wraps them into the ecc types
func GenerateKey() (PublicKey, SecretKey, error) {
	pk, sk, err := ed25519.GenerateKey(rand.Reader)
	return PublicKey(pk), SecretKey(sk), err
}

// Generates keys and writes them to stdout, file or both
// Path parameters can be left if not writing to a file
// `generateKeys` should be one of "stdout", "file" or "both"
// Warning: this function can call Fatal as it is meant to be run as a cli user util.
func KeyGenerationUtil(generateKeys string, pkPath string, skPath string) {
	pk, sk, err := GenerateKey()
	if err != nil {
		log.Fatal("Key generation failed: " + err.Error())
		return
	}

	if generateKeys == generateKeysSTDOUT || generateKeys == generateKeysBoth {
		pkPEM, err := pk.ToPem()
		if err != nil {
			log.Fatal("Couldn't generate pem: " + err.Error())
		}
		skPEM, err := sk.ToPem()
		if err != nil {
			log.Fatal("Couldn't generate pem: " + err.Error())
		}

		// Marshalling to json is used to easily make the keys safe to print
		skPEMstring, err := json.Marshal(skPEM)
		if err != nil {
			log.Fatal("Couldn't marshal: " + err.Error())
		}
		pkPEMMarshall, err := json.Marshal(pkPEM)
		if err != nil {
			log.Fatal("Couldn't marshal: " + err.Error())
		}

		fmt.Printf("PUBLIC_KEY=%s\n", pkPEMMarshall)
		fmt.Printf("SECRET_KEY=%s\n", skPEMstring)
	}

	if generateKeys == generateKeysFile || generateKeys == generateKeysBoth {
		if err := pk.WriteToFile(pkPath); err != nil {
			log.Fatalf("Writing to file failed: %s", err.Error())
		}

		if err := sk.WriteToFile(skPath); err != nil {
			log.Fatalf("Writing to file failed: %s", err.Error())
		}
	}
}
