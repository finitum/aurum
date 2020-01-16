package passwords

import (
	"github.com/trustelem/zxcvbn"
)

var commonPasswords = []string{
	"aurum",
	"finitum",
}

func VerifyPassword(password string, userinput []string) bool {
	if len(password) < 8 {
		return false
	}

	// Max length for bcrypt
	if len(password) > 72 {
		return false
	}

	disallowed := append(commonPasswords, userinput...)
	res := zxcvbn.PasswordStrength(password, disallowed)

	return res.Score > 2
}
