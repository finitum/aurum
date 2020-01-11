package hash

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("yeet")
	assert.Nil(t, err)

	assert.True(t, CheckPasswordHash("yeet", hash))
}
