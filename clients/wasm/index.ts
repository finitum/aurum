import {err, ok, Result} from "neverthrow";

export interface Claims {
    Username: string,
    Refresh: boolean,
    aud: string
    exp: number,
    jti: string,
    iat: number,
    iss: string,
    nbf: number,
    sub: string,
}

interface Error {
    error: string
}

import "./wasm_exec"

export function verifyJwt(token: string, pem: string): Result<Claims, string> {
    // @ts-ignore
    const res = window.ZZZ_AurumWasm_VerifyToken(token, pem)
    if (res.error) {
        return err(res.error)
    }

    return ok(res)
}
