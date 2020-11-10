import {Role} from "./User";

interface IClaimJSON {
    username: string;
    role: Role;
    refresh: boolean;
    exp: number; // Expiry
    iat: number; // Issued At
    nbf: number; // Not Before
}

export interface ITokenPairJSON {
    login_token?: string;
    refresh_token?: string;
}

export class Claims {
    constructor(public username: string,
                public role: number,
                public refresh: boolean,
                public expiresAt: Date,
                public issuedAt: Date,
                public notBefore: Date) { }

    static cachedClaims: [Claims, string] | null = null;

    static fromObject({username, role, refresh, exp, iat, nbf}: IClaimJSON): Claims {
        return new Claims(username, role, refresh, new Date(exp * 1000), new Date(iat * 1000), new Date(nbf * 1000));
    }

    static parse(token: string): Claims {
        if (this.cachedClaims != null && this.cachedClaims[1] == token) {
            return this.cachedClaims[0];
        }

        const part = token.split(".")[1];
        const string = atob(part);

        const res = Claims.fromObject(JSON.parse(string));

        this.cachedClaims = [res, token];
        return res;
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
    static fromObject(object: ITokenPairJSON): TokenPair {

        return new TokenPair(
              object.login_token ? object.login_token : "",
            object.refresh_token ? object.refresh_token : "");
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
    
    json(): string {
        return JSON.stringify({"login_token": this.loginToken, "refresh_token": this.refreshToken});
    }
}
