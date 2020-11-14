import {Client} from "./Client";
import {LocalstorageAvailable} from "./LocalstorageAvailable";
import {ApplicationWithRole, AurumError, ErrorCode, TokenPair, User} from "./Models";
import {err, ok, Result} from "neverthrow";

const LocalStorageRefreshTokenKey = "REFRESH_TOKEN";
const LocalStorageLoginTokenKey = "LOGIN_TOKEN";


export class SingleUserClient {
    private tokenpair: TokenPair | null;
    private user: User | null;
    private client: Client;
    private unAuthorizedHandlers: ((err: AurumError | null) => void)[];

    constructor(baseurl: string) {
        this.tokenpair = null;
        this.user = null;
        this.client = new Client(baseurl)

        this.unAuthorizedHandlers = []

        if (LocalstorageAvailable()) {
            const refresh = localStorage.getItem(LocalStorageRefreshTokenKey)
            const login = localStorage.getItem(LocalStorageLoginTokenKey)

            if (refresh !== null && login !== null) {
                this.tokenpair = {
                    refresh_token: refresh,
                    login_token: login,
                }
            }
        }
    }

    AddUnauthorizedHandler(handler: (err: AurumError | null) => void) {
        this.unAuthorizedHandlers.push(handler)
    }

    async Login(user: User): Promise<Result<null, AurumError>> {
        const tp = await this.client.Login(user);

        if (tp.isOk()) {
            this.tokenpair = tp.value

            if (LocalstorageAvailable() && this.tokenpair !== null) {
                localStorage.setItem(LocalStorageLoginTokenKey, this.tokenpair.login_token)
                localStorage.setItem(LocalStorageRefreshTokenKey, this.tokenpair.refresh_token)
            }

            return ok(null);
        } else {
            return err(tp.error);
        }
    }

    async Register(user: User): Promise<Result<null, AurumError>> {
        const error = await this.client.Register(user);
        if (error.isErr()) {
            return err(error.error)
        }

        return await this.Login(user);
    }

    async GetUserInfo(): Promise<Result<User, AurumError>> {
        if (this.user !== null) {
            return ok(this.user);
        }

        const user = await this.RetryOnUnauthorized((v) => this.client.GetUserInfo(v));

        if (user.isOk()) {
            this.user = user.value;
        }

        return user
    }

    async UpdateUser(user: User): Promise<Result<User, AurumError>> {
        const new_user = await this.RetryOnUnauthorized((v) => this.client.UpdateUser(v, user));

        if (new_user.isOk()) {
            this.user = new_user.value;
        }

        return new_user;
    }

    private async RetryOnUnauthorized<T>(func: (tp: TokenPair, ...args: any[]) => Promise<Result<T, AurumError>>, ...args: any[]): Promise<Result<T, AurumError>> {
        if (this.tokenpair == null) {
            return err({
                Message: "No token stored",
                Code: ErrorCode.Unauthorized,
            })
        }

        const first = await func(this.tokenpair, ...args);

        if (first.isErr() && first.error.Code === ErrorCode.Unauthorized) {
            const new_tp = await this.client.Refresh(this.tokenpair)
            if (new_tp.isOk()) {
                this.tokenpair = new_tp.value
            }

            if (LocalstorageAvailable() && this.tokenpair !== null) {
                localStorage.setItem(LocalStorageLoginTokenKey, this.tokenpair.login_token)
            }

            const second = await func(this.tokenpair, ...args);

            if (second.isErr() && second.error.Code == ErrorCode.Unauthorized) {
                for (const h of this.unAuthorizedHandlers) {
                    h(second.error as AurumError)
                }
            }

            return second;
        } else {
            return first
        }
    }

    async GetApplicationsForUser(user?: User): Promise<Result<ApplicationWithRole[], AurumError>> {
        let checkedUser: User;

        if (typeof user === "undefined") {
            const userOrErr = await this.GetUserInfo()
            if (userOrErr.isOk()) {
                checkedUser = userOrErr.value;
            } else {
                return err(userOrErr.error);
            }
        } else {
            checkedUser = user
        }

        return await this.RetryOnUnauthorized((v) => this.client.GetApplicationsForUser(v, checkedUser));
    }

    Logout() {
        this.tokenpair = null;
        this.user = null;

        if (LocalstorageAvailable()) {
            localStorage.removeItem(LocalStorageLoginTokenKey)
            localStorage.removeItem(LocalStorageRefreshTokenKey)
        }

        for (const h of this.unAuthorizedHandlers) {
            h(null)
        }
    }

    IsLoggedIn(): boolean {
        return this.tokenpair !== null;
    }

    private LoginToken(): string | null {
        if (this.tokenpair === null) {
            return null
        }
        return this.tokenpair.login_token
    }
}