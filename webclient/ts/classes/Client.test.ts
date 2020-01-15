import Client, {ErrorState} from "./Client";

import {generateValidJWT} from "../test-helpers/helpers";
import User from "./User";

import "jest-fetch-mock";
import {TokenPair} from "./jwt";

describe("#Client", () => {
    it("Should create an API and State on construction", () => {
        const b = "https://example.com";
        const c = new Client(b);
        expect(c.api.baseURL).toBe(b);
        expect(c.state).not.toBeNull();
    });
});
describe("#Client.checkLogin()", () => {
    it("Should return a invalidCredentials if no token is available is set", () => {
       const c = new Client("");
       return c.checkLogin().then((r) => {
           expect(r).toStrictEqual([null, ErrorState.InvalidCredentials]);
       });
    });

    it("Should request userinfo on checklogin", () => {
        const c = new Client("");
        const tp = new TokenPair(generateValidJWT(), generateValidJWT(true));
        jest.spyOn(c.state, "tokenPair", "get").mockReturnValue(tp);

        const u = new User("victor", "");
        fetchMock.mockResponseOnce(JSON.stringify(u));

        return c.checkLogin().then(([ur, err]) => {
            expect(ur).toStrictEqual(u);
            expect(err).toStrictEqual(ErrorState.Ok);

            expect(fetchMock.mock.calls.length).toEqual(1);
            expect(fetchMock.mock.calls[0][0]).toEqual("/me");
            expect(fetchMock.mock.calls[0][1]).toEqual({
                method: "GET",
                headers: {
                    "Authorization": `Bearer ${tp.loginToken}`
                },
            });
        });
    });

    it("Should handle server error", () => {
        const c = new Client("");
        const tp = new TokenPair(generateValidJWT(), generateValidJWT(true));
        jest.spyOn(c.state, "tokenPair", "get").mockReturnValue(tp);

        fetchMock.mockResponseOnce("Oopsie Woopsie", {status: 500});

        return c.checkLogin().then(([ur, err]) => {
            expect(ur).toBeNull();
            expect(err).toStrictEqual(ErrorState.ServerError);
        });
    });

    it("Should handle invalid credentials", () => {
        const c = new Client("");
        const tp = new TokenPair(generateValidJWT(), generateValidJWT(true));
        jest.spyOn(c.state, "tokenPair", "get").mockReturnValue(tp);

        fetchMock.mockResponseOnce("Invalid credentials", {status: 401});

        return c.checkLogin().then(([ur, err]) => {
            expect(ur).toBeNull();
            expect(err).toStrictEqual(ErrorState.InvalidCredentials);
        });
    });

    it("Should handle generic error", () => {
        const c = new Client("");
        const tp = new TokenPair(generateValidJWT(), generateValidJWT(true));
        jest.spyOn(c.state, "tokenPair", "get").mockReturnValue(tp);

        fetchMock.mockResponseOnce("smh", {status: 400});

        return c.checkLogin().then(([ur, err]) => {
            expect(ur).toBeNull();
            expect(err).toStrictEqual(ErrorState.Other);
        });
    });

    // it("Should refresh on an expired token", () => {
    //     const c = new Client("");
    //     const tp = new TokenPair(generateExpiredJWT(), generateValidJWT(true));
    //     jest.spyOn(c.state, "tokenPair", "get").mockReturnValue(tp);
    //
    //     // const u = new User("victor", "");
    //     // fetchMock.mockResponseOnce(JSON.stringify(u));
    //
    //     return c.checkLogin().then(([ur, err]) => {
    //
    //     });
    // });
});