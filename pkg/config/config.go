package config

import (
	"github.com/finitum/aurum/pkg/jwt/ecc"
	log "github.com/sirupsen/logrus"
	"go.deanishe.net/env"
)

// TODO: Add options to configure the database
// A struct containing the various config options of Aurum
type EnvConfig struct {
	WebAddr  string `env:"WEB_ADDRESS"`
	BasePath string `env:"BASE_PATH"`

	NoKeyGen   bool `env:"NO_KEY_GENERATE"`
	NoKeyWrite bool `env:"NO_KEY_WRITE"`

	PublicKey string `env:"PUBLIC_KEY"`
	SecretKey string `env:"SECRET_KEY"`

	PublicKeyPath string `env:"PUBLIC_KEY_PATH"`
	SecretKeyPath string `env:"SECRET_KEY_PATH"`

	DgraphUrl string `env:"DGRAPH_URL"`
}

type Config struct {
	WebAddr  string
	BasePath string

	PublicKey ecc.PublicKey
	SecretKey ecc.SecretKey

	DgraphUrl string
}

func defaultEnvConfig() EnvConfig {
	return EnvConfig{
		WebAddr:       "0.0.0.0:8042",
		BasePath:      "/",
		PublicKeyPath: "./id_25519.pub",
		SecretKeyPath: "./id_25519",
		SecretKey:     "",
		PublicKey:     "",
		NoKeyGen:      false,
		NoKeyWrite:    false,
		DgraphUrl:     "localhost:9080",
	}
}

func GetEnvConfig(e ...env.Env) *EnvConfig {
	dc := defaultEnvConfig()

	if err := env.Bind(&dc, e...); err != nil {
		log.Fatal(err.Error())
	}

	return &dc
}

// Helper function for getting the default config values
func GetConfig(e ...env.Env) *Config {
	ec := GetEnvConfig(e...)
	pk, sk, err := findKeys(ec)
	if err != nil {
		log.Fatal(err.Error())
	}

	return &Config{
		WebAddr:   ec.WebAddr,
		BasePath:  ec.BasePath,
		PublicKey: pk,
		SecretKey: sk,

		DgraphUrl: ec.DgraphUrl,
	}
}

// EphemeralConfig returns a config that doesn't read from env or write to files, mainly for use in tests
func EphemeralConfig() *Config {
	ec := defaultEnvConfig()
	ec.NoKeyWrite = true

	pk, sk, err := findKeys(&ec)
	if err != nil {
		log.Fatal(err.Error())
	}

	return &Config{
		WebAddr:   ec.WebAddr,
		BasePath:  ec.BasePath,
		PublicKey: pk,
		SecretKey: sk,
	}
}
