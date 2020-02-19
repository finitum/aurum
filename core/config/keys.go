package config

import (
	"aurum/jwt/ecc"
	log "github.com/sirupsen/logrus"
	"os"
)

// TODO: Testing
func (b *Builder) FindKeys(noWrite bool) BuilderProcess {

	// If Public key in env
	if publickey := os.Getenv("PUBLIC_KEY"); publickey != "" {
		pk, err := ecc.FromPem([]byte(publickey))
		if err != nil {
			log.Panic("Could not parse public key given in environment variable.")
		}
		b.PublicKey = pk.(ecc.PublicKey)
		log.Info("Using public key from environment")
	} else {
		// Else read from file
		pk, err := ecc.FromFile(b.PubKeyPath)
		if err != nil {
			log.Info("Could not find public key in file or environment. Generating from secret key if it exists.")
		} else {
			b.PublicKey = pk.(ecc.PublicKey)
		}
	}

	// Read secret key from env
	if secretkey := os.Getenv("SECRET_KEY"); secretkey != "" {
		sk, err := ecc.FromPem([]byte(secretkey))
		if err != nil {
			log.Panic("Could not parse public key given in environment variable.")
		}
		b.SecretKey = sk.(ecc.SecretKey)
		log.Info("Using private key from environment")
	} else {
		// Else read from file

		sk, err := ecc.FromFile(b.SecretKeyPath)
		if err != nil {
			log.Warn("Couldn't find keys in environment, and reading the secret key failed. Generating secret key. " +
				"NOTE: this is normal on the first run of aurum.")
		} else {
			// found sk file path
			b.SecretKey = sk.(ecc.SecretKey)
		}
	}

	// If key is nil -> generate
	if b.SecretKey == nil && b.NoKeyGen {
		log.Panic("Secret key is nil and not allowed to generate")
	} else if b.SecretKey == nil && !b.NoKeyGen {
		pk, sk, err := ecc.GenerateKey()
		if err != nil {
			log.Panic("Key generation went wrong: " + err.Error())
		}
		b.PublicKey = pk
		b.SecretKey = sk

		if !noWrite {
			if err := b.PublicKey.WriteToFile(b.PubKeyPath); err != nil {
				log.Fatalf("Could't write private key to file: %s", err)
			}

			if err := b.SecretKey.WriteToFile(b.SecretKeyPath); err != nil {
				log.Fatalf("Could't write public key to file: %s", err)
			}
		}
	}

	return b
}
