import User, {IUser, Role} from "./User";

describe("#User", () => {
    it("Should be able to create a User from an object", () => {
        const user = {
            username: "victor",
            password: "pass",
            email: "mail",
            role: Role.Admin,
            blocked: false
        };

        const created = User.fromObject(user);
        expect(created).toEqual(user);
    });

    it("Should have fallbacks for undefined fields", () => {
        const user = {
            username: "victor",
        };

        const expected: IUser = {
            username: "victor",
            password: "",
            email: "",
            role: Role.User,
            blocked: false
        };

        const created = User.fromObject(user as IUser);
        expect(created).toEqual(expected);
    });
});
