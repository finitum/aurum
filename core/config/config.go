package config

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"os"
)

// A struct containing the various config options of Aurum
type Config struct {
	JWTKey  []byte
	WebAddr string
	Path    string
}

// An interface for Config builder
type BuilderProcess interface {
	// Sets config options to their default
	SetDefault() BuilderProcess
	// Gets config options from env vars
	SetFromEnvironment() BuilderProcess
	// Gets config options from a specified file
	SetFromFile(path string) BuilderProcess
	// Builds the config
	Build() *Config
}

// A builder to build the config
type Builder struct {
	Config
}

func (b *Builder) SetDefault() BuilderProcess {
	b.Config = Config{
		JWTKey:  []byte("ChangeMe"),
		WebAddr: "127.0.0.1:8042",
		Path:    "/",
	}
	return b
}

func (b *Builder) SetFromEnvironment() BuilderProcess {
	if jwt := os.Getenv("JWTKEY"); jwt != "" {
		b.JWTKey = []byte(jwt)
	}

	if web := os.Getenv("WEBURL"); web != "" {
		b.WebAddr = web
	}

	if path := os.Getenv("PATH"); path != "" {
		b.Path = path
	}

	return b
}

func (b *Builder) SetFromFile(path string) BuilderProcess {
	panic("implement me")
}

func (b *Builder) Build() *Config {
	if bytes.Equal(b.JWTKey, []byte("ChangeMe")) {
		log.Warn("The JWTKey has not been changed from the default")
	}

	return &b.Config
}

// Helper function for getting the default config values
func GetDefault() *Config {
	return new(Builder).SetDefault().Build()
}
