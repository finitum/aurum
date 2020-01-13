package passwords

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerifyPasswordTooShort(t *testing.T) {
	assert.False(t, VerifyPassword("18828"));
}

func TestVerifyPasswordTooLong(t *testing.T) {
	assert.False(t, VerifyPassword("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"));
}

func TestVerifyPasswordTooCommon(t *testing.T) {
	assert.False(t, VerifyPassword("password"));
}