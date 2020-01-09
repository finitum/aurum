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
		JWTKey: []byte("ChangeMe"),
		WebURL: "127.0.0.1:8042",
	}, b)
}

func TestBuilder_SetFromEnvironment(t *testing.T) {
	bd := new(Builder)

	_ = os.Setenv("JWTKEY", "key")
	_ = os.Setenv("WEBURL", "url")
	b := bd.SetFromEnvironment().Build()
	_ = os.Setenv("JWTKEY", "")
	_ = os.Setenv("WEBURL", "")

	assert.Equal(t, &Config{
		JWTKey: []byte("key"),
		WebURL: "url",
	}, b)
}

func TestBuilder_SetDefault_SetFromEnvironment(t *testing.T) {
	bd := new(Builder)
	_ = os.Setenv("JWTKEY", "key")
	b := bd.SetDefault().SetFromEnvironment().Build()
	_ = os.Setenv("JWTKEY", "")

	assert.Equal(t, &Config{
		JWTKey: []byte("key"),
		WebURL: "127.0.0.1:8042",
	}, b)
}