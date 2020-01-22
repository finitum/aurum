import Client, {API, ErrorState} from "./Client";

import {generateValidJWT} from "./__TEST__/helpers";
import User from "./User";

import "jest-fetch-mock";
import {TokenPair} from "./jwt";



describe("#API", () => {

    describe("login", () => {

        it("Should be able to login successfully", () => {
            const a = new API("");
            const u = new User("victor", "");
            const tp = new TokenPair(generateValidJWT(), generateValidJWT(true));
            fetchMock.mockResponseOnce(tp.json());

            return a.login(u).then(([tpr, err]) => {
                expect(tpr).toStrictEqual(tp);
                expect(err).toStrictEqual(ErrorState.Ok);

                expect(fetchMock.mock.calls.length).toEqual(1);
                expect(fetchMock.mock.calls[0][0]).toEqual("/login");
                expect(fetchMock.mock.calls[0][1]).toEqual({
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify(u)
                });
            });
        });

        it("Should be able to handle invalid credentials", () => {
            const a = new API("");
            const u = new User("victor", "");

            fetchMock.mockResponseOnce("Invalid Credentials", {status: 401});

            return a.login(u).then(([tpr, err]) => {
                expect(tpr).toBeNull();
                expect(err).toStrictEqual(ErrorState.InvalidCredentials);

                expect(fetchMock.mock.calls.length).toEqual(1);
                expect(fetchMock.mock.calls[0][0]).toEqual("/login");
                expect(fetchMock.mock.calls[0][1]).toEqual({
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify(u)
                });
            });
        });

        it("Should be able to handle generic server error", () => {
            const a = new API("");
            const u = new User("victor", "");

            fetchMock.mockResponseOnce("Internal Server Error", {status: 500});

            return a.login(u).then(([tpr, err]) => {
                expect(tpr).toBeNull();
                expect(err).toStrictEqual(ErrorState.ServerError);

                expect(fetchMock.mock.calls.length).toEqual(1);
                expect(fetchMock.mock.calls[0][0]).toEqual("/login");
                expect(fetchMock.mock.calls[0][1]).toEqual({
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify(u)
                });
            });
        });

        it("Should be able to handle unexpected errors", () => {
            const a = new API("");
            const u = new User("victor", "");

            fetchMock.mockResponseOnce("Not enough wooloos", {status: 420});

            return a.login(u).then(([tpr, err]) => {
                expect(tpr).toBeNull();
                expect(err).toStrictEqual(ErrorState.Other);

                expect(fetchMock.mock.calls.length).toEqual(1);
                expect(fetchMock.mock.calls[0][0]).toEqual("/login");
                expect(fetchMock.mock.calls[0][1]).toEqual({
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify(u)
                });
            });
        });
    });

    describe("signup", () => {
        it("Should be able to register successfully", () => {
            const a = new API("");
            const u = new User("victor", "");
            fetchMock.mockResponseOnce("", {status: 201});

            return a.signup(u).then((err) => {
                expect(err).toStrictEqual(ErrorState.Ok);

                expect(fetchMock.mock.calls.length).toEqual(1);
                expect(fetchMock.mock.calls[0][0]).toEqual("/signup");
                expect(fetchMock.mock.calls[0][1]).toEqual({
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify(u)
                });
            });
        });

        it("Should be able to handle conflict error", () => {
            const a = new API("");
            const u = new User("victor", "");
            fetchMock.mockResponseOnce("User already exists", {status: 409});

            return a.signup(u).then((err) => {
                expect(err).toStrictEqual(ErrorState.UserExists);

                expect(fetchMock.mock.calls.length).toEqual(1);
                expect(fetchMock.mock.calls[0][0]).toEqual("/signup");
                expect(fetchMock.mock.calls[0][1]).toEqual({
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify(u)
                });
            });
        });

        it("Should be able to handle InvalidPassword error", () => {
            const a = new API("");
            const u = new User("victor", "");
            fetchMock.mockResponseOnce("Password not sufficient", {status: 422});

            return a.signup(u).then((err) => {
                expect(err).toStrictEqual(ErrorState.InvalidPasswordError);

                expect(fetchMock.mock.calls.length).toEqual(1);
                expect(fetchMock.mock.calls[0][0]).toEqual("/signup");
                expect(fetchMock.mock.calls[0][1]).toEqual({
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify(u)
                });
            });
        });

        it("Should be able to handle generic error", () => {
            const a = new API("");
            const u = new User("victor", "");
            fetchMock.mockResponseOnce("wooloo", {status: 500});

            return a.signup(u).then((err) => {
                expect(err).toStrictEqual(ErrorState.ServerError);

                expect(fetchMock.mock.calls.length).toEqual(1);
                expect(fetchMock.mock.calls[0][0]).toEqual("/signup");
                expect(fetchMock.mock.calls[0][1]).toEqual({
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify(u)
                });
            });
        });

        it("Should be able to handle unexpected errors", () => {
            const a = new API("");
            const u = new User("victor", "");
            fetchMock.mockResponseOnce("wooloo?", {status: 420});

            return a.signup(u).then((err) => {
                expect(err).toStrictEqual(ErrorState.Other);

                expect(fetchMock.mock.calls.length).toEqual(1);
                expect(fetchMock.mock.calls[0][0]).toEqual("/signup");
                expect(fetchMock.mock.calls[0][1]).toEqual({
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify(u)
                });
            });
        });
    });

    describe("getMe", () => {

        it("Should be able to get the user successfully", () => {
            const a = new API("");
            const token = generateValidJWT();
            const u = new User("victor", "");

            fetchMock.mockResponseOnce(JSON.stringify(u));

            return a.getMe(token).then(([ru, err]) => {
                expect(err).toStrictEqual(ErrorState.Ok);
                expect(ru).toStrictEqual(u);

                expect(fetchMock.mock.calls.length).toEqual(1);
                expect(fetchMock.mock.calls[0][0]).toEqual("/user");
                expect(fetchMock.mock.calls[0][1]).toEqual({
                    method: "GET",
                    headers: {
                        "Authorization": `Bearer ${token}`
                    }
                });
            });
        });

    });
});


describe("#Client", () => {
    it("Should create an API and State on construction", () => {
        const b = "https://example.com";
        const c = new Client(b);
        expect(c.api.baseURL).toBe(b);
        expect(c.state).not.toBeNull();
    });

    describe("checkLogin()", () => {
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
                expect(fetchMock.mock.calls[0][0]).toEqual("/user");
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
    });
});