import {ITokenPairJSON, TokenPair} from "./jwt";
import User from "./User";
import State from "../State/State";
import {domstate} from "../globals";
import {DOMState} from "../DOM/DOMStateManager";
import Config from "../Config";

export enum ErrorState {
    Ok,
    InvalidCredentials,
    UserExists,
    ServerError,
    InvalidPasswordError,

    Other,
}

/**
 * @private
 * The class which is used to communicate to the backend,
 * it wraps all HTTP methods into JS functions for ease of use.
 * please use the [Client] class whenever necessary.
 */
export class API {
    baseURL: string;

    constructor(baseURL: string) {
        this.baseURL = baseURL;
    }

    /**
     * Logs a user in to the application.
     * @param user for authentication, username and password need to be non-null
     * @returns [TokenPair] or [ErrorState] depending on if it went correctly
     */
    async login(user: User): Promise<[TokenPair | null, ErrorState]> {
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

        if (res.status == 422) {
            return ErrorState.InvalidPasswordError;
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
    async getMe(token: string): Promise<[User | null, ErrorState]> {
        const res = await fetch(`${this.baseURL}/user`,{
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
    async refresh(tokenPair: TokenPair): Promise<[TokenPair | null, ErrorState]> {
        const res = await fetch(`${this.baseURL}/refresh`, {
            method: "POST",
            body: tokenPair.json()
        });

        if (res.status.toString().startsWith("5")) {
            return [null, ErrorState.ServerError];
        }
        if (res.status == 401) {
            return [null, ErrorState.InvalidCredentials];
        }

        if (!res.status.toString().startsWith("2")) {
            return [null, ErrorState.InvalidCredentials];
        }

        const resultObject: ITokenPairJSON = await res.json();
        if (resultObject.login_token === undefined) {
            return [null, ErrorState.ServerError];
        }

        const newTokenPair = new TokenPair(resultObject.login_token, tokenPair.refreshToken);

        return [newTokenPair, ErrorState.Ok];
    }

    /**
     * Changes the password of a user.
     * @param user The user with a new password set
     * @param token the token used for authentication
     */
    async updateUser(token: string, user: User): Promise<ErrorState> {
        const res = await fetch(`${this.baseURL}/user`, {
            method: "PUT",
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
    async getUsers(token: string, start: number, end: number): Promise<[User[] ,ErrorState]> {
        const res = await fetch(`${this.baseURL}/users`, {
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
    public user: User | null = null;
    public api: API;
    public state: State;

    private readonly worker: Worker | null;

    private static instance: Client | null = null;

    constructor(baseURL: string) {
        this.api = new API(baseURL);
        this.state = new State();

        //check if webworkers are available
        if(window.Worker){
            this.worker = new Worker("../Worker/worker.ts");
            this.worker.onmessage = async (e: MessageEvent): Promise<void> => this.onWorkerChange(e);
        } else {
            this.worker = null;
        }
    }

    public static getInstance(): Client {
        if(Client.instance == null){
            Client.instance = new Client(Config.API_URL);
        }

        return Client.instance;
    }

    // Saves the tokenpair in the state and notifies the webworker
    private set tokenPair(tp: TokenPair) {
        this.state.tokenPair = tp;

        if(this.worker != null){
            this.worker.postMessage(this.state.tokenPair);
        }
    }

    // When we receive a message from a worker,
    // set the received tokenpair to be the localstorage stored one.
    private async onWorkerChange(e: MessageEvent): Promise<void> {

        if (this.state.tokenPair != null && this.state.tokenPair.isRefreshValid) {
            const tp = await this.api.refresh(this.state.tokenPair);

            if (tp === null) {
                domstate.change(DOMState.Login);
            } else if(tp[0] != null) {
                this.tokenPair = tp[0];
            } else {
                domstate.change(DOMState.Login);
            }
        } else {
            domstate.change(DOMState.Login);
        }
    }

    /**
     * Checks localstorage for authentication tokesns and automatically
     * logs the user in if the tokens are still valid.
     */
    async checkLogin(): Promise<[User | null, ErrorState]> {
        let tokenPair, err, newuser;

        if (this.state.tokenPair) {
            if(this.state.tokenPair.isLoginValid){
                [newuser, err] = await this.api.getMe(this.state.tokenPair.loginToken);
                this.user = newuser;

                this.tokenPair = this.state.tokenPair;

                return [this.user, err];
            } else if (this.state.tokenPair.isRefreshValid) {
                // refresh
                const tp = await this.api.refresh(this.state.tokenPair);

                // if the refresh failed, just log in again.
                if(tp === null) {
                    return [null, ErrorState.InvalidCredentials];
                }

                [tokenPair, err]  = await this.api.refresh(this.state.tokenPair);
                if (err != ErrorState.Ok || tokenPair == null) {
                    return [null, err];
                }

                [newuser, err] = await this.api.getMe(tokenPair.loginToken);
                if (err != ErrorState.Ok) {
                    return [null, err];
                }

                this.user = newuser;

                this.tokenPair = tokenPair;

                return [this.user, ErrorState.Ok];
            }
        }
        return [null, ErrorState.InvalidCredentials];
    }

    /**
     * Logs a user in
     * @param username the user to log in
     * @param password the password for said user
     */
    async login(username: string, password: string): Promise<[User | null, ErrorState]> {

        const [tokenPair, err1] = await this.api.login(new User(username, password));

        if(err1 !== ErrorState.Ok || tokenPair == null) {
            return [null, err1];
        }

        const [newuser, err2] = await this.api.getMe(tokenPair.loginToken);
        if(err2 !== ErrorState.Ok) {
            return [null, err2];
        }

        this.user = newuser;

        this.tokenPair = tokenPair;

        return [this.user, ErrorState.Ok];
    }

    /**
     * Creates a user
     * @param username The username of the user
     * @param password The password of the user
     * @param email The email of the user
     * @returns A [User] object on success or an [ErrorState] on failure
     */
    async signup(username: string, password: string, email: string): Promise<[User | null, ErrorState]> {

        const err = await this.api.signup(new User(username, password, email));
        if (err !== ErrorState.Ok) {
            return [null, err];
        }

        return await this.login(username, password);
    }

    /**
     * Removes the tokens from localstorage
     */
    logout(): void {
        this.state.tokenPair = null;

        if (this.worker !== null) {
            this.worker.postMessage(null);
        }
    }

    /**
     * Changes the password of a user
     * @param password the new password
     * @returns the new [User] object or an [ErrorState]
     */
    async changePassword(password: string): Promise<[User | null, ErrorState]> {
        if (this.user === null || this.state.tokenPair == null) {
            return [null, ErrorState.InvalidCredentials];
        }

        const newUser = this.user;
        newUser.password = password;

        const err = await this.api.updateUser(this.state.tokenPair.loginToken, newUser);
        if (err !== ErrorState.Ok) {
            return [null, err];
        }

        const [getNewUser, err2] = await this.api.getMe(this.state.tokenPair.loginToken);
        if(err2 != ErrorState.Ok || getNewUser == null) {
            return [null, err2];
        }

        this.user = getNewUser;

        return [this.user, ErrorState.Ok];
    }

    /**
     * Gets all users with pagination
     * @param start of the page
     * @param end of the page
     * @returns a [User] Array or an [ErrorState]
     */
    async getUsers(start: number, end: number): Promise<[User[], ErrorState]> {
        if(this.state.tokenPair == null) {
            return [[], ErrorState.InvalidCredentials];
        }

        return this.api.getUsers(this.state.tokenPair.loginToken, start, end);
    }

    /**
     * Blocks a user
     *
     * @param user The user to be blocked, with the blocker set according
     * @returns an ErrorState
     */
    async setBlocked(user: User): Promise<ErrorState> {
        if(this.state.tokenPair == null) {
            return ErrorState.InvalidCredentials;
        }

        // `==` to also check for undef
        if(user.blocked == null || user.username == null) {
            return ErrorState.Other;
        }

        return this.api.updateUser(this.state.tokenPair.loginToken, user);
    }
}