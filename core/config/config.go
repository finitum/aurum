package config

import (
	"aurum/jwt/algorithms"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)


// TODO: Add options to configure the database
// A struct containing the various config options of Aurum
type Config struct {
	WebAddr   string
	BasePath  string
	KeyPath   string

	NoKeyGen bool
	SecretKey ed25519.PrivateKey
	PublicKey ed25519.PublicKey
}


// An interface for Config builder
type BuilderProcess interface {
	// Sets config options to their default
	SetDefault() BuilderProcess
	// Gets config options from env vars
	SetFromEnvironment() BuilderProcess
	// Gets the keys
	FindKeys(bool) BuilderProcess
	// Builds the config
	Build() *Config
}

// A builder to build the config
type Builder struct {
	Config
}

func (b *Builder) SetDefault() BuilderProcess {
	b.Config = Config{
		WebAddr:  "0.0.0.0:8042",
		BasePath: "/",
		KeyPath: "/tmp/aurum-keys",
		SecretKey: nil,
		PublicKey: nil,
		NoKeyGen: false,
	}
	return b
}

func (b *Builder) SetFromEnvironment() BuilderProcess {

	if webaddr := os.Getenv("WEBADDR"); webaddr != "" {
		b.WebAddr = webaddr
	}

	if basepath := os.Getenv("BASEPATH"); basepath != "" {
		b.BasePath = basepath
	}

	if keypath := os.Getenv("KEYPATH"); keypath != "" {
		b.KeyPath = keypath
	}

	if nokeygen := os.Getenv("NOKEYGEN"); nokeygen != "" {
		b.NoKeyGen = true
	}

	return b
}


const publicKeyPemHeader = "PUBLIC KEY"
const ed25519PrivateKeyPemHeader = "ED25519 PRIVATE KEY" // TODO: Not sure if this is the correct label

func (b *Builder) FindKeys(noWrite bool) BuilderProcess {

	// If Public key in env
	if publickey := os.Getenv("PUBLIC_KEY"); publickey != "" {
		pk, err := algorithms.DecodeEdDSAPublicKey([]byte(publickey))
		if err != nil {
			log.Panic("Could not parse public key given in environment variable.")
		}
		b.PublicKey = pk
		log.Info("Using public key from environment")
	} else {
		// Else read from file

		pkeyfile, err := ioutil.ReadFile(b.KeyPath + "/key_pub")
		if err != nil {
			log.Info("Could not find public key in file or environment. Generating from private key if it exists.")
		} else {
			pk, err := algorithms.DecodeEdDSAPublicKey(pkeyfile)
			if err != nil {
				log.Panic("Could not parse public key given in file variable.")
			}
			b.PublicKey = pk
		}
	}

	// Read secret key from env
	if secretkey := os.Getenv("SECRET_KEY"); secretkey != "" {
		sk, err := algorithms.DecodeEdDSAPrivatekey([]byte(secretkey))
		if err != nil {
			log.Panic("Could not parse public key given in environment variable.")
		}
		b.SecretKey = sk
		log.Info("Using private key from environment")
	} else {
		// Else read from file

		skeyfile, err := ioutil.ReadFile(b.KeyPath + "/key")
		if err != nil {
			log.Warn("Couldn't find keys in environment, and reading the secret key failed. Generating private key. " +
				"NOTE: this is normal on the first run of aurum.")
		} else {
			// found sk file path
			sk, err := algorithms.DecodeEdDSAPrivatekey(skeyfile)
			if err != nil {
				log.Panic("Could not parse private key given in file variable.")
			}
			b.SecretKey = sk
		}
	}

	// If key is nil -> generate
	if b.SecretKey == nil && b.NoKeyGen {
		log.Panic("Secret key is nil and not allowed to generate")
	} else if b.SecretKey == nil && !b.NoKeyGen {
		pk, sk, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			log.Panic("Key generation went wrong: " + err.Error())
		}
		b.PublicKey = pk
		b.SecretKey = sk

		if !noWrite {
			if err := b.writeKeyToFile(sk, ed25519PrivateKeyPemHeader); err != nil {
				log.Fatalf("Could't write private key to file: %s", err)
			}

			if err := b.writeKeyToFile(pk, publicKeyPemHeader); err != nil {
				log.Fatalf("Could't write public key to file: %s", err)
			}
		}
	}

	return b
}

func (b *Builder) writeKeyToFile(key []byte, keytype string) error {
	block := pem.Block{
		Type:   keytype,
		Bytes:  key[:],
	}

	var fileLoc string

	switch keytype {
		case publicKeyPemHeader:
			fileLoc = b.KeyPath + "/id_ed25519.pub"
			break
		case ed25519PrivateKeyPemHeader:
			fileLoc = b.KeyPath + "/id_ed25519"
			break
		default:
			return errors.New("unrecognized key type")
	}

	writer, err := os.Open(fileLoc)
	if err != nil {
		return err
	}
	defer writer.Close()

	err = pem.Encode(writer, &block)
	if err != nil {
		return err
	}

	return nil
}

func (b *Builder) Build() *Config {
	return &b.Config
}

// Helper function for getting the default config values
func GetDefault() *Config {
	return new(Builder).SetDefault().FindKeys(true).Build()
}
