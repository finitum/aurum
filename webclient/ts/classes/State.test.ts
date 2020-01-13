import State, {Storage} from "./State";
import {TokenPair} from "./jwt";


describe("#Storage", () => {
   it("Should parse from JSON correctly", () => {
       const tp =new TokenPair("a", "b");

       const s = {
            tokenPair: tp
       };

       const json = JSON.stringify(s);

       const res = Storage.fromJSON(json);

       expect(res.tokenPair).toEqual(tp);
   });
   it("Should parse to JSON correctly", () => {
        const tp =new TokenPair("a", "b");

        const s = {
            tokenPair: tp
        };

        const json = JSON.stringify(s);

        const res: Storage = Storage.fromJSON(json);

        expect(res.json()).toBe(json);
    });
});

describe("#State", () => {

    // Clear localstorage before all State tests
    beforeEach(() => {
        localStorage.clear();
    });

    it("Should try to get from localstorage on construction", () => {
        new State();
        expect(localStorage.getItem).toHaveBeenCalledWith("storage");
    });

    it("Should store to localstorage on construction", () => {
        const s = new State();
        expect(localStorage.getItem).toHaveBeenCalledWith("storage");

        s.tokenPair = new TokenPair("a", "b");
        expect(localStorage.setItem).toHaveBeenCalled();

        const expected = {
            tokenPair: {
                loginToken: "a",
                refreshToken: "b"
            }
        };

        const json = JSON.stringify(expected);

        expect(localStorage.__STORE__["storage"]).toBe(json);
        expect(localStorage.length).toBe(1);
    });

    it("Should return null when no token is stored", () => {
        const s = new State();
        expect(localStorage.length).toBe(0);
        expect(s.tokenPair).toBeNull();
    });

    it("Should return a stored tokenpair", () => {
        // Setup localstorage
        localStorage.__STORE__["storage"] = JSON.stringify({
            tokenPair: {
                loginToken: "a",
                refreshToken: "b"
            }
        });

        const s = new State();
        expect(localStorage.getItem).toHaveBeenCalledWith("storage");
        expect(s.tokenPair.loginToken).toBe("a");
        expect(s.tokenPair.refreshToken).toBe("b");
    });
});