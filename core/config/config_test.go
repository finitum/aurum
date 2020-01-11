package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestBuilder_SetDefault(t *testing.T) {
	bd := new(Builder)
	b := bd.SetDefault().Build()

	assert.Equal(t, &Config{
		JWTKey:  []byte("ChangeMe"),
		WebAddr: "127.0.0.1:8042",
		Path:    "/",
	}, b)
}

func TestGetDefault(t *testing.T) {
	bd := new(Builder).SetDefault().Build()
	assert.Equal(t, bd, GetDefault())
}

func TestBuilder_SetFromEnvironment(t *testing.T) {
	bd := new(Builder)

	_ = os.Setenv("JWTKEY", "key")
	_ = os.Setenv("WEBURL", "url")
	_ = os.Setenv("PATH", "/asd")
	b := bd.SetFromEnvironment().Build()
	_ = os.Setenv("JWTKEY", "")
	_ = os.Setenv("WEBURL", "")
	_ = os.Setenv("PATH", "")

	assert.Equal(t, &Config{
		JWTKey:  []byte("key"),
		WebAddr: "url",
		Path:    "/asd",
	}, b)
}

func TestBuilder_SetDefault_SetFromEnvironment(t *testing.T) {
	bd := new(Builder)
	_ = os.Setenv("JWTKEY", "key")
	b := bd.SetDefault().SetFromEnvironment().Build()
	_ = os.Setenv("JWTKEY", "")

	assert.Equal(t, &Config{
		JWTKey:  []byte("key"),
		WebAddr: "127.0.0.1:8042",
		Path:    "/",
	}, b)
}
