package main

import (
	"github.com/finitum/aurum/pkg/jwt"
	"github.com/finitum/aurum/pkg/jwt/ecc"
	"syscall/js"
)

//go:generate sh -c "cp $(go env GOROOT)/misc/wasm/wasm_exec.js ."

func main() {
	js.Global().Set("ZZZ_AurumWasm_VerifyToken", VerifyTokenWrapper())
	select {}
}

// Signature is (token, pem) -> claims
func VerifyTokenWrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 2 {
			return MarshalError("VerifyToken: expected two arguments")
		}

		// Arg 0 == token
		token := args[0].String()

		// Arg 1 == pem
		pem := args[1].String()
		key, err := ecc.FromPem([]byte(pem))
		if err != nil {
			return MarshalError("VerifyToken: could not decode pem: " + err.Error())
		}

		pk := key.(ecc.PublicKey)

		// Call function
		claims, err := jwt.VerifyJWT(token, pk)
		if err != nil {
			return MarshalError("VerifyToken: could not decode pem: " + err.Error())
		}

		if err := claims.Valid(); err != nil {
			return MarshalError("VerifyToken: invalid claims: " + err.Error())
		}

		// Return object
		return MarshalClaims(claims)
	})
}

func MarshalError(err string) js.Value {
	obj := make(map[string]interface{}, 1)
	obj["error"] = err
	return js.ValueOf(obj)
}

func MarshalClaims(claims *jwt.Claims) js.Value {
	obj := make(map[string]interface{}, 9)

	obj["Username"] = claims.Username
	obj["Refresh"] = claims.Refresh

	obj["aud"] = claims.Audience
	obj["exp"] = claims.ExpiresAt
	obj["jti"] = claims.Id
	obj["iat"] = claims.IssuedAt
	obj["iss"] = claims.Issuer
	obj["nbf"] = claims.NotBefore
	obj["sub"] = claims.Subject

	return js.ValueOf(obj)
}
