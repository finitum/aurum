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

/**
 * The class which is used to communicate to the backend,
 * it wraps all HTTP methods into JS functions for ease of use.
 */
class API {
    baseURL: string;

    constructor(baseURL: string) {
        this.baseURL = baseURL;
    }

    /**
     * Logs a user in to the application.
     * @param user for authentication, username and password need to be non-null
     * @returns [TokenPair] or [ErrorState] depending on if it went correctly
     */
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

    /**
     * Registers a user
     * @param user the user to register
     * @returns An [ErrorState] to signify the result
     */
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

    /**
     * Gets the user object from the server of the currently logged in user
     * @param token the token to use for authentication
     */
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

    /**
     * Uses the refresh token to get a new login token
     * @param tokenPair the tokenpair containg the refresh token used for getting the new login token
     * @returns A new tokenpair or an [ErrorState]
     */
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

    /**
     * Changes the password of a user.
     * @param user The user with a new password set
     * @param token the token used for authentication
     */
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

    /**
     * Gets all users (paginated)
     * For this function the user needs to be Admin.
     * @param start of the page
     * @param end end of the page
     * @param token authentication token.
     */
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

/**
 * This class abstracts the [API] class to make interacting with the backend easier.
 */
export default class Client {
    public user: User;
    public api: API;
    public state: State;

    constructor(baseURL: string) {
        this.api = new API(baseURL);
        this.state = new State();
    }

    /**
     * Checks localstorage for authentication tokesns and automatically
     * logs the user in if the tokens are still valid.
     */
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

    /**
     * Logs a user in
     * @param username the user to log in
     * @param password the password for said user
     */
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

    /**
     * Creates a user
     * @param username The username of the user
     * @param password The password of the user
     * @param email The email of the user
     * @returns A [User] object on success or an [ErrorState] on failure
     */
    async signup(username: string, password: string, email: string): Promise<[User, ErrorState]> {

        const err = await this.api.signup(new User(username, password, email));
        if (err !== ErrorState.Ok) {
            return [null, err];
        }

        return await this.login(username, password)
    }

    /**
     * Removes the tokens from localstorage
     */
    logout(): void {
        this.state.tokenPair = null;
    }

    /**
     * Changes the password of a user
     * @param password the new password
     * @returns the new [User] object or an [ErrorState]
     */
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

    /**
     * Gets all users with pagination
     * @param start of the page
     * @param end of the page
     * @returns a [User] Array or an [ErrorState]
     */
    async getUsers(start: number, end: number): Promise<[User[], ErrorState]> {
        return this.api.getUsers(start, end, this.state.tokenPair.loginToken);
    }
}