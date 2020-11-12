import axios, {AxiosInstance} from "axios";
import {AurumError, TokenPair, User} from "./Models";
import {LocalstorageAvailable} from "./LocalstorageAvailable";

const LocalStorageRefreshTokenKey = "REFRESH_TOKEN";
const LocalStorageLoginTokenKey = "LOGIN_TOKEN";

export class Client {
    private axios: AxiosInstance;
    private user: User | null;
    private tokenpair: TokenPair | null;

    constructor(baseurl: string) {
        this.axios = axios.create({
            baseURL: baseurl,
            headers: {
                "Content-type": "application/json"
            }
        });

        this.axios.interceptors.request.use(conf => {
            if (this.LoginToken() != null) {
                conf.headers.Authorization = `Bearer ${this.LoginToken()}`;
            }

            return conf;
        })

        this.user = null;
        this.tokenpair = null;

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

    // Login makes a request to log in a user. It returns either an error or null.
    // In the null case login was succesful and the Aurum client succesfully saved the token.
    async Login(user: User): Promise<AurumError | null> {
        const resp = await this.axios.post("/login", user)

        if (resp.status === 200) {
            this.tokenpair = resp.data as TokenPair;

            if (LocalstorageAvailable()) {
                localStorage.setItem(LocalStorageLoginTokenKey, this.tokenpair.login_token)
                if(this.tokenpair.refresh_token !== null) {
                    localStorage.setItem(LocalStorageRefreshTokenKey, this.tokenpair.refresh_token)
                }
            }

            return null
        } else {
            return resp.data as AurumError
        }
    }

    async Register(user: User) {

        // TODO: Handle response & also login
    }

    async GetUser(): Promise<User | AurumError> {
        if (this.user !== null) {
            return this.user
        }

        const resp = await this.axios.get("/user")

        if (resp.status === 200) {
            return resp.data as AurumError
        } else {
            return resp.data as AurumError
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