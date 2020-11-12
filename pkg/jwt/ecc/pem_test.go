package ecc

import (
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/assert"
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

func TestStuff(t *testing.T) {
	pk, sk, err := GenerateKey()
	assert.NoError(t, err)

	pkB64 := base64.StdEncoding.EncodeToString(pk)
	skB64 := base64.StdEncoding.EncodeToString(sk)

	pkPem, err := pk.ToPem()
	assert.NoError(t, err)

	fmt.Printf("\nPublic Key PEM:\n%v", pkPem)
	fmt.Printf("\nPublic Key Base64: %v", pkB64)
	fmt.Printf("\nSecret Key Base64: %v\n", skB64)
}
