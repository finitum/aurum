import Client, {API, ErrorState} from "./Client";

import {generateExpiredJWT, generateValidJWT} from "./__TEST__/helpers";
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

    describe("refresh", () => {
        it("should refresh correctly", () => {
            const a = new API("");

            const tp = new TokenPair("", generateValidJWT());
            const responseTP = new TokenPair(generateValidJWT(), "");

            fetchMock.mockResponseOnce(responseTP.json());

            return a.refresh(tp).then( ([ntp, err]) => {
                expect(err).toStrictEqual(ErrorState.Ok);

                const resultTP = new TokenPair(responseTP.loginToken, tp.refreshToken);
                expect(ntp).toEqual(resultTP);

                expect(fetchMock.mock.calls.length).toEqual(1);
                expect(fetchMock.mock.calls[0][0]).toEqual("/refresh");

                expect(fetchMock.mock.calls[0][1]).toEqual({
                   method: "POST",
                   body: tp.json(),
                });
            });
        });

        it("should catch server error", () => {
            const a = new API("");

            const tp = new TokenPair("", generateValidJWT());
            const responseTP = new TokenPair(generateValidJWT(), "");

            fetchMock.mockResponseOnce("wooloo", {status: 500});

            return a.refresh(tp).then( ([ntp, err]) => {
                expect(err).toStrictEqual(ErrorState.ServerError);
                expect(ntp).toBeNull();
            });
        });

        it("should handle invalid creds", () => {
            const a = new API("");

            const tp = new TokenPair("", "");

            fetchMock.mockResponseOnce("wooloo", {status: 401});

            return a.refresh(tp).then( ([ntp, err]) => {
                expect(err).toStrictEqual(ErrorState.InvalidCredentials);
                expect(ntp).toBeNull();
            });
        });

        it("should handle other errors", () => {
            const a = new API("");

            const tp = new TokenPair("", "");

            fetchMock.mockResponseOnce("wooloo", {status: 444});

            return a.refresh(tp).then( ([ntp, err]) => {
                expect(err).toStrictEqual(ErrorState.InvalidCredentials);
                expect(ntp).toBeNull();
            });
        });

        it("should handle invalid response", () => {
            const a = new API("");

            const tp = new TokenPair("", "");

            fetchMock.mockResponseOnce("{}");

            return a.refresh(tp).then( ([ntp, err]) => {
                expect(err).toStrictEqual(ErrorState.ServerError);
                expect(ntp).toBeNull();
            });
        });
    });

    describe("updateUser", () => {
        it("should update correctly", () => {
            const a = new API("");
            const u = new User("victor", "");

            const token = generateValidJWT();

            fetchMock.mockResponseOnce("", {status: 200});


            return a.updateUser(token, u).then( (err) => {
                expect(err).toStrictEqual(ErrorState.Ok);

                expect(fetchMock.mock.calls.length).toEqual(1);
                expect(fetchMock.mock.calls[0][0]).toEqual("/user");

                expect(fetchMock.mock.calls[0][1]).toEqual({
                    method: "PUT",
                    body: JSON.stringify(u),
                    headers: {
                        "Authorization": `Bearer ${token}`
                    },
                });
            });
        });

        it("should catch server error", () => {
            const a = new API("");
            const u = new User("victor", "");

            const token = generateValidJWT();

            fetchMock.mockResponseOnce("", {status: 500});

            return a.updateUser(token, u).then( (err) => {
                expect(err).toStrictEqual(ErrorState.ServerError);
            });
        });

        it("should handle invalid creds", () => {
            const a = new API("");
            const u = new User("victor", "");
            const token = generateValidJWT();


            fetchMock.mockResponseOnce("wooloo", {status: 401});

            return a.updateUser(token, u).then( (err) => {
                expect(err).toStrictEqual(ErrorState.InvalidCredentials);
            });
        });

        it("should handle other errors", () => {
            const a = new API("");
            const u = new User("victor", "");
            const token = generateValidJWT();


            fetchMock.mockResponseOnce("wooloo", {status: 444});

            return a.updateUser(token, u).then( (err) => {
                expect(err).toStrictEqual(ErrorState.Other);
            });
        });
    });

    describe("getUsers", () => {
        it("should get correctly", () => {
            const a = new API("");

            const u1 = new User("v", "");
            const u2 = new User("j", "");
            const u = [u1, u2];

            const token = generateValidJWT();

            fetchMock.mockResponseOnce(JSON.stringify(u));


            return a.getUsers(token, 0, 1).then( ([users, err]) => {
                expect(err).toStrictEqual(ErrorState.Ok);
                expect(users).toEqual(u);

                expect(fetchMock.mock.calls.length).toEqual(1);
                expect(fetchMock.mock.calls[0][0]).toEqual("/users?start=0&end=1");

                expect(fetchMock.mock.calls[0][1]).toEqual({
                    method: "GET",
                    headers: {
                        "Authorization": `Bearer ${token}`
                    },
                });
            });
        });

        it("should catch server error", () => {
            const a = new API("");

            const token = generateValidJWT();

            fetchMock.mockResponseOnce("", {status: 500});

            return a.getUsers(token, 0, 1).then( ([u,err]) => {
                expect(u).toStrictEqual([]);
                expect(err).toStrictEqual(ErrorState.ServerError);
            });
        });

        it("should handle invalid creds", () => {
            const a = new API("");
            const token = generateValidJWT();


            fetchMock.mockResponseOnce("wooloo", {status: 401});

            return a.getUsers(token, 0, 1).then( ([u,err]) => {
                expect(u).toStrictEqual([]);
                expect(err).toStrictEqual(ErrorState.InvalidCredentials);
            });
        });

        it("should handle other errors", () => {
            const a = new API("");
            const token = generateValidJWT();

            fetchMock.mockResponseOnce("wooloo", {status: 444});

            return a.getUsers(token, 0, 1).then( ([u,err]) => {
                expect(u).toStrictEqual([]);
                expect(err).toStrictEqual(ErrorState.Other);
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

        it("Should refresh the login token if refresh is valid but login isn't", () => {
            const c = new Client("");
            const tp = new TokenPair(generateExpiredJWT(), generateValidJWT(true));
            const newtp = new TokenPair(generateValidJWT(), generateValidJWT(true));
            const u = new User("a", "");
            jest.spyOn(c.state, "tokenPair", "get").mockReturnValue(tp);
            // @ts-ignore we actually use a higher js version than intellij thinks
            jest.spyOn(c.api, "refresh").mockReturnValue(new Promise(resolve => {
                resolve([newtp, ErrorState.Ok]);
            }));
            // @ts-ignore we actually use a higher js version than intellij thinks
            jest.spyOn(c.api, "getMe").mockReturnValue(new Promise(resolve => {
                resolve([u, ErrorState.Ok]);
            }));

            return c.checkLogin().then(([ru,err]) => {
                expect(err).toStrictEqual(ErrorState.Ok);
                expect(ru).toBe(u);
            });
        });

        it("Should handle invalid tp response on refresh", () => {
            const c = new Client("");
            const tp = new TokenPair(generateExpiredJWT(), generateValidJWT(true));
            jest.spyOn(c.state, "tokenPair", "get").mockReturnValue(tp);
            // @ts-ignore we actually use a higher js version than intellij thinks
            jest.spyOn(c.api, "refresh").mockReturnValue(new Promise(resolve => {
                resolve([null, ErrorState.InvalidCredentials]);
            }));

            return c.checkLogin().then(([ru,err]) => {
                expect(err).toStrictEqual(ErrorState.InvalidCredentials);
                expect(ru).toBeNull();
            });
        });

        it("Should handle invalid getMe response while refresh", () => {
            const c = new Client("");
            const tp = new TokenPair(generateExpiredJWT(), generateValidJWT(true));
            const newtp = new TokenPair(generateValidJWT(), generateValidJWT(true));
            jest.spyOn(c.state, "tokenPair", "get").mockReturnValue(tp);
            // @ts-ignore we actually use a higher js version than intellij thinks
            jest.spyOn(c.api, "refresh").mockReturnValue(new Promise(resolve => {
                resolve([newtp, ErrorState.Ok]);
            }));
            // @ts-ignore we actually use a higher js version than intellij thinks
            jest.spyOn(c.api, "getMe").mockReturnValue(new Promise(resolve => {
                resolve([null, ErrorState.Other]);
            }));

            return c.checkLogin().then(([ru,err]) => {
                expect(err).toStrictEqual(ErrorState.Other);
                expect(ru).toBeNull();
            });
        });
    });


});