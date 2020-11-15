package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/finitum/aurum/pkg/jwt"
	"github.com/finitum/aurum/pkg/jwt/ecc"
	"syscall/js"
)

//go:generate sh -c "cp $(go env GOROOT)/misc/wasm/wasm_exec.js ."

func Warn(msg ...interface{}) {
	console := js.Global().Get("console")
	warn := console.Get("warn")
	warn.Invoke(msg...)
}

func Warnf(format string, args ...interface{}) {
	str := fmt.Sprintf(format, args...)
	Warn(str)
}

func main() {
	Warnf("Initialized Go wasm lib: %v", errors.New("some error"))

	js.Global().Set("VerifyToken", VerifyTokenWrapper())
	<-make(chan struct{})
}

// Signature is (token, pem) -> claims
func VerifyTokenWrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 2 {
			Warn("VerifyToken: expected two arguments")
			return ""
		}

		// Arg 0 == token
		token := args[0].String()

		// Arg 1 == pem
		pem := args[1].String()
		key, err := ecc.FromPem([]byte(pem))
		if err != nil {
			Warnf("VerifyToken: could not decode pem: %v", err)
			return ""
		}

		pk := key.(ecc.PublicKey)

		// Call function
		claims, err := jwt.VerifyJWT(token, pk)
		if err != nil {
			Warnf("VerifyToken: could not decode pem: %v", err)
			return ""
		}

		if err := claims.Valid(); err != nil {
			Warnf("VerifyToken: invalid claims: %v", err)
			return ""
		}

		// Return json
		ret, err := json.Marshal(claims)
		if err != nil {
			Warnf("VerifyToken: couldn't marshal claims to json: %v", err)
			return ""
		}

		return string(ret)
	})
}
