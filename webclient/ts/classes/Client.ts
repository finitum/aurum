import {TokenPair} from "./jwt";
import User from "./User";
import State from "./State";

export enum ErrorState {
    Ok,
    InvalidCredentials,
    UserExists,
    ServerError,

    Other,
}

class API {
    baseURL: string;

    constructor(baseURL: string) {
        this.baseURL = baseURL;
    }

    async login(user: User): Promise<[TokenPair, ErrorState]> {
        const res = await fetch(`${this.baseURL}/login`,{
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(user),
        });

        if (res.status == 401) {
            return [null, ErrorState.InvalidCredentials];
        }

        if (res.status.toString().startsWith("5")) {
            return [null, ErrorState.ServerError];
        }

        if (!res.status.toString().startsWith("2")) {
            return [null, ErrorState.Other];
        }

        return [TokenPair.fromObject(await res.json()), ErrorState.Ok];
    }

    async signup(user: User): Promise<ErrorState> {

        const res = await fetch(`${this.baseURL}/signup`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(user),
        });


        if (res.status == 409) {
            return ErrorState.UserExists;
        }

        if (res.status.toString().startsWith("5")) {
            return ErrorState.ServerError;
        }

        if (!res.status.toString().startsWith("2")) {
            return ErrorState.Other;
        }

        return ErrorState.Ok;
    }

    async getMe(token: string): Promise<[User, ErrorState]> {
        const res = await fetch(`${this.baseURL}/me`,{
            method: "GET",
            headers: {
                "Authorization": `Bearer ${token}`
            },
        });

        if (res.status.toString().startsWith("5")) {
            return [null, ErrorState.ServerError];
        }

        if (res.status == 401) {
            return [null, ErrorState.InvalidCredentials];
        }

        if (!res.status.toString().startsWith("2")) {
            return [null, ErrorState.Other];
        }


        return [User.fromObject(await res.json()), ErrorState.Ok];
    }

    async refresh(tokenPair: TokenPair): Promise<[TokenPair, ErrorState]> {
        const res = await fetch(`${this.baseURL}/refresh`, {
            method: "POST",
            body: JSON.stringify(tokenPair)
        });

        if (res.status.toString().startsWith("5")) {
            return [null, ErrorState.ServerError];
        }

        if (res.status == 401) {
            return [null, ErrorState.InvalidCredentials];
        }

        if (!res.status.toString().startsWith("2")) {
            return null;
        }

        return [TokenPair.fromJSON(await res.json()), ErrorState.Ok];
    }

    // accepts a user with a changed password
    async changePassword(user: User, token: string): Promise<ErrorState> {
        const res = await fetch(`${this.baseURL}/changepassword`, {
            method: "POST",
            body: JSON.stringify(user),
            headers: {
                "Authorization": `Bearer ${token}`
            },
        });

        if (res.status.toString().startsWith("5")) {
            return ErrorState.ServerError;
        }

        if (res.status == 401) {
            return ErrorState.InvalidCredentials;
        }

        if (!res.status.toString().startsWith("2")) {
            return ErrorState.Other;
        }

        return ErrorState.Ok;
    }

    async getUsers(start: number, end: number, token: string): Promise<[User[] ,ErrorState]> {
        const res = await fetch(`${this.baseURL}/requestusers`, {
            method: "POST",
            body: JSON.stringify({start: start, end: end}),
            headers: {
                "Authorization": `Bearer ${token}`
            },
        });

        if (res.status.toString().startsWith("5")) {
            return [[], ErrorState.ServerError];
        }

        if (res.status == 401) {
            return [[], ErrorState.InvalidCredentials];
        }

        if (!res.status.toString().startsWith("2")) {
            return [[], ErrorState.Other];
        }

        return [await res.json(), ErrorState.Ok];
    }
}

export default class Client {
    public user: User;
    public api: API;
    public state: State;

    constructor(baseURL: string) {
        this.api = new API(baseURL);
        this.state = new State();
    }

    async checkLogin(): Promise<[User, ErrorState]> {
        let tokenPair, err, newuser;

        if (this.state.tokenPair) {
            if(this.state.tokenPair.isLoginValid()){
                [newuser, err] = await this.api.getMe(this.state.tokenPair.loginToken);
                this.user = newuser;

                return [this.user, ErrorState.Ok];
            } else if (this.state.tokenPair.isRefreshValid) {
                // refresh
                const tp = await this.api.refresh(this.state.tokenPair);

                // if the refresh failed, just log in again.
                if(tp === null) {
                    return [null, ErrorState.InvalidCredentials];
                }

                [tokenPair, err]  = await this.api.refresh(this.state.tokenPair);
                if (err != ErrorState.Ok) {
                    return [null, err];
                }

                [newuser, err] = await this.api.getMe(tokenPair.loginToken);
                if (err != ErrorState.Ok) {
                    return [null, err];
                }

                this.user = newuser;
                this.state.tokenPair = tokenPair;


                return [this.user, ErrorState.Ok];
            }
        } else {
            return [null, ErrorState.InvalidCredentials];
        }
    }

    async login(username: string, password: string): Promise<[User, ErrorState]> {

        const [tokenPair, err1] = await this.api.login(new User(username, password));

        if(err1 !== ErrorState.Ok) {
            return [null, err1];
        }

        const [newuser, err2] = await this.api.getMe(tokenPair.loginToken);
        if(err2 !== ErrorState.Ok) {
            return [null, err2];
        }

        this.user = newuser;
        this.state.tokenPair = tokenPair;

        return [this.user, ErrorState.Ok];
    }



    async signup(username: string, password: string, email: string): Promise<[User, ErrorState]> {

        const err = await this.api.signup(new User(username, password, email));
        if (err !== ErrorState.Ok) {
            return [null, err];
        }

        const [tokenPair, err2] = await this.api.login(new User(username, password));

        if(err2 !== ErrorState.Ok) {
            return [null, err2];
        }

        const [newuser, err3] = await this.api.getMe(tokenPair.loginToken);
        if(err3 !== ErrorState.Ok) {
            return [null, err3];
        }

        this.user = newuser;
        this.state.tokenPair = tokenPair;

        return [this.user, ErrorState.Ok];
    }

    logout(): void {
        this.state.tokenPair = null;
    }

    async changePassword(password: string): Promise<[User, ErrorState]> {
        if (this.user === null) {
            return [null, ErrorState.InvalidCredentials];
        }

        const err = await this.api.changePassword(new User("", password), this.state.tokenPair.loginToken);
        if (err !== ErrorState.Ok) {
            return [null, err];
        }

        return [this.user, ErrorState.Ok];
    }

    async getUsers(start: number, end: number): Promise<[User[], ErrorState]> {
        return this.api.getUsers(start, end, this.state.tokenPair.loginToken);
    }
}