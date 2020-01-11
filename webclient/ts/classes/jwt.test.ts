import {Claims, isJWTValid} from "./jwt";

describe("#isJWTValid()", () => {
    it("Should fail on invalid token", () => {
        expect(isJWTValid("a.e30=")).toBe(false);
    });
    it("Should succeed on a valid token", () => {
        const now = new Date();
        now.setMinutes(now.getMinutes() - 1);

        const future = new Date();
        future.setHours(future.getHours() + 1);

        const fakeTokenClaims = {
            username: "victor",
            role: 0,
            refresh: false,
            exp: future.getTime() / 1000,
            iat: now.getTime() / 1000,
            nbf: now.getTime() / 1000
        };

        const base64 = btoa(JSON.stringify(fakeTokenClaims));

        expect(isJWTValid(`a.${base64}`)).toBe(true);
    });
    it("Should fail on an expired token", () => {
        const iat = new Date();
        iat.setHours(iat.getHours() - 2);

        const exp = new Date();
        exp.setHours(exp.getHours() - 1);

        const fakeTokenClaims = {
            username: "victor",
            role: 0,
            refresh: false,
            exp: exp.getTime() / 1000,
            iat: iat.getTime() / 1000,
            nbf: iat.getTime() / 1000
        };

        const base64 = btoa(JSON.stringify(fakeTokenClaims));

        expect(isJWTValid(`a.${base64}`)).toBe(false);
    });
    it("should fail on a not yet valid toekn", () => {
        const iat = new Date();
        iat.setHours(iat.getHours() + 2);

        const exp = new Date();
        exp.setHours(exp.getHours() + 3);

        const fakeTokenClaims = {
            username: "victor",
            role: 0,
            refresh: false,
            exp: exp,
            iat: iat,
            nbf: iat
        };

        const base64 = btoa(JSON.stringify(fakeTokenClaims));

        expect(isJWTValid(`a.${base64}`)).toBe(false);
    });
});

describe("Claims", () => {
    it("Should parse object correctly", () => {
        const now = new Date();

        const claimsObj = {
            username: "victor",
            role: 0,
            refresh: false,
            exp: now.getTime() / 1000,
            iat: now.getTime() / 1000,
            nbf: now.getTime() / 1000
        };

        const c = Claims.fromObject(claimsObj);

        expect(c.username).toBe(claimsObj.username);
        expect(c.role).toBe(claimsObj.role);
        expect(c.refresh).toBe(claimsObj.refresh);
        expect(c.expiresAt).toStrictEqual(now);
        expect(c.issuedAt).toStrictEqual(now);
        expect(c.notBefore).toStrictEqual(now);
    });

    it("Should parse JWT string correctly", () => {
        const token = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE1Nzg2NzgxODcsImV4cCI6MTU3ODY3OTM5MiwiYXVkIjoiIiwic3ViIjoiIiwidXNlcm5hbWUiOiJ2aWN0b3IifQ.6SQZiU_wjZSwZ3315t27aDhGzbevBq1mj3EPdv2960s";

        const c = Claims.parse(token);
        expect(c.username).toBe("victor");

        expect(c.issuedAt).toStrictEqual(new Date(1578678187 * 1000));
        expect(c.expiresAt).toStrictEqual(new Date(1578679392 * 1000));
    });

    it("Should parse JSONified object correctly", () => {
        const now = new Date();

        const claimsObj = {
            username: "victor",
            role: 0,
            refresh: false,
            exp: now.getTime() / 1000,
            iat: now.getTime() / 1000,
            nbf: now.getTime() / 1000
        };

        const json = JSON.stringify(claimsObj);
        const b64 = btoa(json);
        const tkn = `a.${b64}.c`;

        const c = Claims.parse(tkn);

        expect(c.username).toBe(claimsObj.username);
        expect(c.role).toBe(claimsObj.role);
        expect(c.refresh).toBe(claimsObj.refresh);
        expect(c.expiresAt).toStrictEqual(now);
    });
});