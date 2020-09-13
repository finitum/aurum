package passwords

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerifyPasswordTooShort(t *testing.T) {
	assert.False(t, VerifyPassword("18828", nil))
}

func TestVerifyPasswordTooLong(t *testing.T) {
	assert.False(t, VerifyPassword("4b93310ed64ce510889be78f32203f9768c4054b9af08489ed90a59465616ef64b93310ed64ce510889be78f32203f9768c4054b9af08489ed90a59465616ef6", nil))
}

func TestVerifyPasswordTooCommon(t *testing.T) {
	assert.False(t, VerifyPassword("password", nil))
}

func TestVerifyPasswordSameAsCommonWords(t *testing.T) {
	assert.False(t, VerifyPassword("finitumaurum", nil))
}

func TestVerifyPasswordSameAsUserInput(t *testing.T) {
	assert.False(t, VerifyPassword("7da033bd32005113f2208eb87bc94c126a42aadf0c94065b1fa4d9d68e7c318f", []string{
		"2aadf0c94065b1fa4d9d68e7c318f",
		"7da033bd32005113f2208eb87bc94c126a4",
	}))
}

func TestValidPassword(t *testing.T) {
	assert.True(t, VerifyPassword("7da033bd32005113f2208eb87bc94c126a42aadf0c94065b1fa4d9d68e7c318f", nil))
}
