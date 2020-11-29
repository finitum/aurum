package passwords

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyPasswordTooShort(t *testing.T) {
	assert.False(t, CheckStrength("18828", nil))
}

func TestVerifyPasswordTooLong(t *testing.T) {
	assert.False(t, CheckStrength("4b93310ed64ce510889be78f32203f9768c4054b9af08489ed90a59465616ef64b93310ed64ce510889be78f32203f9768c4054b9af08489ed90a59465616ef6", nil))
}

func TestVerifyPasswordTooCommon(t *testing.T) {
	assert.False(t, CheckStrength("password", nil))
}

func TestVerifyPasswordSameAsCommonWords(t *testing.T) {
	assert.False(t, CheckStrength("finitumaurum", nil))
}

func TestVerifyPasswordSameAsUserInput(t *testing.T) {
	assert.False(t, CheckStrength("7da033bd32005113f2208eb87bc94c126a42aadf0c94065b1fa4d9d68e7c318f", []string{
		"2aadf0c94065b1fa4d9d68e7c318f",
		"7da033bd32005113f2208eb87bc94c126a4",
	}))
}

func TestValidPassword(t *testing.T) {
	assert.True(t, CheckStrength("7da033bd32005113f2208eb87bc94c126a42aadf0c94065b1fa4d9d68e7c318f", nil))
}
