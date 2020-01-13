package passwords

import (
	log "github.com/sirupsen/logrus"
	"github.com/trustelem/zxcvbn"
)

var commonPasswords = []string {
	"aurum",
	"finitum",
}

func VerifyPassword(password string, userinput []string) bool {
	disallowed := append(commonPasswords, userinput...)
	res := zxcvbn.PasswordStrength(password, disallowed)

	log.Trace(res)

	return res.Score > 2
}
