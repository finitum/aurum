package config

import (
	"github.com/finitum/aurum/internal/jwt/ecc"
	"github.com/test-go/testify/assert"
	"os"
	"testing"
)

const nonExistingPath = "./non-existing-path/"

func assertNonExistingPathDoesntExist(t *testing.T) {
	if _, err := os.Stat(nonExistingPath); !os.IsNotExist(err) {
		t.Fail()
	}
}

func TestInvalidPath(t *testing.T) {
	assertNonExistingPathDoesntExist(t)

	_, _, err := findKeys(&EnvConfig{
		WebAddr:       "",
		BasePath:      "",
		NoKeyGen:      true,
		NoKeyWrite:    true,
		PublicKey:     "",
		SecretKey:     "",
		PublicKeyPath: nonExistingPath,
		SecretKeyPath: "",
	})

	assert.Error(t, err)

	_, _, err = findKeys(&EnvConfig{
		WebAddr:       "",
		BasePath:      "",
		NoKeyGen:      true,
		NoKeyWrite:    true,
		PublicKey:     "",
		SecretKey:     "",
		PublicKeyPath: "",
		SecretKeyPath: nonExistingPath,
	})

	assert.Error(t, err)
}

func TestInvalidKey(t *testing.T) {
	_, _, err := findKeys(&EnvConfig{
		WebAddr:       "",
		BasePath:      "",
		NoKeyGen:      true,
		NoKeyWrite:    true,
		PublicKey:     "invalid public key",
		SecretKey:     "",
		PublicKeyPath: "",
		SecretKeyPath: "",
	})

	assert.Error(t, err)

	_, _, err = findKeys(&EnvConfig{
		WebAddr:       "",
		BasePath:      "",
		NoKeyGen:      true,
		NoKeyWrite:    true,
		PublicKey:     "",
		SecretKey:     "invalid secret key",
		PublicKeyPath: "",
		SecretKeyPath: "",
	})

	assert.Error(t, err)
}

func TestNonMatchingKeys(t *testing.T) {
	pk1, _, err := ecc.GenerateKey()
	assert.NoError(t, err)
	_, sk2, err := ecc.GenerateKey()
	assert.NoError(t, err)

	spem, err := sk2.ToPem()
	assert.NoError(t, err)
	ppem, err := pk1.ToPem()
	assert.NoError(t, err)

	_, _, err = findKeys(&EnvConfig{
		WebAddr:       "",
		BasePath:      "",
		NoKeyGen:      true,
		NoKeyWrite:    true,
		PublicKey:     ppem,
		SecretKey:     spem,
		PublicKeyPath: "",
		SecretKeyPath: "",
	})

	assert.Error(t, err)
}

func TestWrongKeys(t *testing.T) {
	pk, sk, err := ecc.GenerateKey()
	assert.NoError(t, err)

	spem, err := sk.ToPem()
	assert.NoError(t, err)
	ppem, err := pk.ToPem()
	assert.NoError(t, err)

	_, _, err = findKeys(&EnvConfig{
		WebAddr:       "",
		BasePath:      "",
		NoKeyGen:      true,
		NoKeyWrite:    true,
		PublicKey:     spem, // Secret key pem passed as public key
		SecretKey:     ppem, // Public key pem passed as secret key
		PublicKeyPath: "",
		SecretKeyPath: "",
	})

	assert.Error(t, err)
}

func TestMatchingKeys(t *testing.T) {
	pk, sk, err := ecc.GenerateKey()
	assert.NoError(t, err)

	spem, err := sk.ToPem()
	assert.NoError(t, err)
	ppem, err := pk.ToPem()
	assert.NoError(t, err)

	npk, nsk, err := findKeys(&EnvConfig{
		WebAddr:       "",
		BasePath:      "",
		NoKeyGen:      true,
		NoKeyWrite:    true,
		PublicKey:     ppem,
		SecretKey:     spem,
		PublicKeyPath: "",
		SecretKeyPath: "",
	})

	assert.NoError(t, err)
	assert.Equal(t, npk, pk)
	assert.Equal(t, nsk, sk)
}

func TestNoSkNoGen(t *testing.T) {
	pk, _, err := ecc.GenerateKey()
	assert.NoError(t, err)

	ppem, err := pk.ToPem()
	assert.NoError(t, err)

	_, _, err = findKeys(&EnvConfig{
		WebAddr:       "",
		BasePath:      "",
		NoKeyGen:      true, // Can't generate
		NoKeyWrite:    true,
		PublicKey:     ppem,
		SecretKey:     "", // No secret key passed
		PublicKeyPath: "",
		SecretKeyPath: "",
	})

	assert.Error(t, err)
}

func TestKeyGen(t *testing.T) {
	pk, sk, err := findKeys(&EnvConfig{
		WebAddr:       "",
		BasePath:      "",
		NoKeyGen:      false, // *Can* generate
		NoKeyWrite:    true,
		PublicKey:     "", // No public key passed
		SecretKey:     "", // No secret key passed
		PublicKeyPath: "",
		SecretKeyPath: "",
	})

	assert.NoError(t, err)

	assert.NotNil(t, pk)
	assert.NotNil(t, sk)

	assert.True(t, sk.Matches(pk))
}

func TestNoKeyGenBecausePublicKey(t *testing.T) {
	pk, _, err := ecc.GenerateKey()
	assert.NoError(t, err)

	ppem, err := pk.ToPem()
	assert.NoError(t, err)

	_, _, err = findKeys(&EnvConfig{
		WebAddr:       "",
		BasePath:      "",
		NoKeyGen:      false, // *Can* generate
		NoKeyWrite:    true,
		PublicKey:     ppem, // A valid public key *is* passed
		SecretKey:     "",   // No secret key passed
		PublicKeyPath: "",
		SecretKeyPath: "",
	})

	assert.Error(t, err)
}

func TestOnlySecretKey(t *testing.T) {
	_, sk, err := ecc.GenerateKey()
	assert.NoError(t, err)
	spem, err := sk.ToPem()
	assert.NoError(t, err)

	npk, nsk, err := findKeys(&EnvConfig{
		WebAddr:       "",
		BasePath:      "",
		NoKeyGen:      false, // *Can* generate
		NoKeyWrite:    true,
		PublicKey:     "",   // No public key is passe
		SecretKey:     spem, // A valid secret key is passed
		PublicKeyPath: "",
		SecretKeyPath: "",
	})

	assert.NoError(t, err)
	assert.Equal(t, sk, nsk)

	assert.NotNil(t, npk)
	assert.NotNil(t, nsk)
	assert.True(t, sk.Matches(npk))
}
