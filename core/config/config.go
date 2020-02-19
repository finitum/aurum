package config

import (
	"aurum/jwt/ecc"
	"github.com/deanishe/go-env"
	"log"
)


// TODO: Add options to configure the database
// A struct containing the various config options of Aurum
type Config struct {
	WebAddr   string				`env:"WEBADDR"`
	BasePath  string				`env:"BASEPATH"`

	NoKeyGen  bool          		`env:"NOKEYGEN"`
	PublicKey ecc.PublicKey 		`env:"PUBLIC_KEY"`
	SecretKey ecc.SecretKey 		`env:"SECRET_KEY"`

	PubKeyPath    string			`env:"PUBLIC_KEY_PATH"`
	SecretKeyPath string			`env:"SECRET_KEY_PATH"`
}


// An interface for Config builder
type BuilderProcess interface {
	// Sets config options to their default
	SetDefault() BuilderProcess
	// Gets config options from env vars
	SetFromEnvironment(...env.Env) BuilderProcess
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
		WebAddr:  		"0.0.0.0:8042",
		BasePath: 		"/",
		PubKeyPath: 	"./id_25519.pub",
		SecretKeyPath:  "./id_25519",
		SecretKey: 		nil,
		PublicKey: 		nil,
		NoKeyGen: 		false,
	}
	return b
}

func (b *Builder) SetFromEnvironment(e ...env.Env) BuilderProcess {
	if err := env.Bind(&b.Config, e...); err != nil {
		log.Fatal(err.Error())
	}

	return b
}

func (b *Builder) Build() *Config {
	return &b.Config
}

// Helper function for getting the default config values
func GetDefault() *Config {
	return new(Builder).SetDefault().FindKeys(true).Build()
}
