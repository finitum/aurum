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

declare global {
    interface Window {
        Go: any
        WebAssembly: any
        ZZZ_AurumWasm_VerifyToken(token: string, pem: string): Claims | Error
    }
}

// @ts-ignore
const g = global || window || self

import "./wasm_exec"

(() => {

    const go = new g.Go();
    g.WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result: any) => {
        go.run(result.instance);
    });
})();

export function verifyJwt(token: string, pem: string): Result<Claims, string> {
    const res = g.ZZZ_AurumWasm_VerifyToken(token, pem)
    if (res.error) {
        return err(res.error)
    }

    return ok(res)
}
