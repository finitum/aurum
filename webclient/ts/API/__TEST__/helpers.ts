export const generateValidJWT = (refresh = false): string => {
    const now = new Date();
    now.setMinutes(now.getMinutes() - 1);

    const future = new Date();
    future.setHours(future.getHours() + 1);

    const fakeTokenClaims = {
        username: "victor",
        role: 0,
        refresh: refresh,
        exp: future.getTime() / 1000,
        iat: now.getTime() / 1000,
        nbf: now.getTime() / 1000
    };

    const base64 = btoa(JSON.stringify(fakeTokenClaims));

    return `a.${base64}.b`;
};

export const generateExpiredJWT = (refresh = false): string => {
    const iat = new Date();
    iat.setHours(iat.getHours() - 2);

    const exp = new Date();
    exp.setHours(exp.getHours() - 1);

    const fakeTokenClaims = {
        username: "victor",
        role: 0,
        refresh: refresh,
        exp: exp.getTime() / 1000,
        iat: iat.getTime() / 1000,
        nbf: iat.getTime() / 1000
    };

    const base64 = btoa(JSON.stringify(fakeTokenClaims));

    return `a.${base64}.b`;
};