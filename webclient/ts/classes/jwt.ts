import {Role} from "./User";

interface IClaimJSON {
    username: string;
    role: Role;
    refresh: boolean;
    exp: number; // Expiry
    iat: number; // Issued At
    nbf: number; // Not Before
}

export class Claims {
    constructor(public username: string,
                public role: number,
                public refresh: boolean,
                public expiresAt: Date,
                public issuedAt: Date,
                public notBefore: Date) { }

    static fromObject({username, role, refresh, exp, iat, nbf}: IClaimJSON): Claims {
        return new Claims(username, role, refresh, new Date(exp * 1000), new Date(iat * 1000), new Date(nbf * 1000));
    }

    static parse(token: string): Claims {
        const part = token.split(".")[1];
        const string = atob(part);

        return Claims.fromObject(JSON.parse(string));
    }
}

export function isJWTValid(tokenstring: string): boolean {
    const claims = Claims.parse(tokenstring);

    const now = new Date();

    const expired = claims.expiresAt < now;
    const nbfPassed = now > claims.notBefore;

    return !expired && nbfPassed;
}

export class TokenPair {

    constructor(public loginToken: string, public refreshToken: string) { }

    // eslint-disable-next-line @typescript-eslint/camelcase
    static fromObject({login_token = "", refresh_token = ""}): TokenPair {
        return new TokenPair(login_token, refresh_token);
    }

    static fromJSON(json: string): TokenPair {
        return TokenPair.fromObject(JSON.parse(json));
    }

    get isLoginValid(): boolean {
        return isJWTValid(this.loginToken);
    }

    get isRefreshValid(): boolean {
        return isJWTValid(this.refreshToken);
    }
}
