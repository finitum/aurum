
export enum Role {
    User = 0,
    Admin = 1,
}

export interface IUser {
    username: string;
    password: string;
    email: string;
    role: Role;
    blocked: boolean;
}

export default class User implements IUser {
    public username: string;
    public password: string;
    public email: string;
    public role: Role;
    public blocked: boolean;

    constructor(username: string, password: string, email = "", role: Role = Role.User, blocked = false) {
        this.username = username;
        this.password = password;
        this.email = email;
        this.role = role;
        this.blocked = blocked;
    }

    static fromObject(x: IUser): User {
        const password = x.password === undefined ? "" : x.password;
        const email = x.email === undefined ? "" : x.email;
        const role = x.role === undefined ? Role.User : x.role;
        const blocked = x.blocked === undefined ? false : x.blocked;
        return new User(x.username, password, email, role, blocked);
    }
}