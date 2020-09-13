package ecc

import (
	"github.com/test-go/testify/assert"
	"testing"
)

func TestPKToFromPem(t *testing.T) {
	pk, _, err := GenerateKey()
	assert.NoError(t, err)

	pkPEM, err := pk.ToPem()
	assert.NoError(t, err)

	pkFromPem, err := FromPem([]byte(pkPEM))
	assert.Equal(t, pk, pkFromPem)
}

func TestSKToFromPem(t *testing.T) {
	_, sk, err := GenerateKey()
	assert.NoError(t, err)

	skPEM, err := sk.ToPem()
	assert.NoError(t, err)

	skFromPem, err := FromPem([]byte(skPEM))
	assert.Equal(t, sk, skFromPem)
}
