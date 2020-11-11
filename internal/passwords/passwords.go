package passwords

import (
	"github.com/trustelem/zxcvbn"
)

var commonTerms = []string{
	"aurum",
	"finitum",
}

func CheckStrength(password string, userinput []string) bool {
	if len(password) < 8 {
		return false
	}

	// Max length for bcrypt
	if len(password) > 72 {
		return false
	}

	disallowed := append(commonTerms, userinput...)
	res := zxcvbn.PasswordStrength(password, disallowed)

	return res.Score > 2
}
