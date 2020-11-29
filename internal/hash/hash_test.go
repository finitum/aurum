package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("yeet")
	assert.Nil(t, err)

	assert.True(t, CheckPasswordHash("yeet", hash))
}
