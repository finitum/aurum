import axios, {AxiosInstance, AxiosRequestConfig} from "axios";
import {ApplicationWithRole, AurumError, ErrorCode, PublicKey, TokenPair, User} from "./Models";
import {err, ok, Result, ResultAsync} from "neverthrow";
import {Claims, verifyJwt} from "aurum-crypto";


export class Client {
    private axios: AxiosInstance;
    private publicKey: string | null;

    constructor(baseurl: string) {
        this.axios = axios.create({
            baseURL: baseurl,
            headers: {
                "Content-type": "application/json"
            }
        });


        this.publicKey = null;

    }

    async GetPublicKey(): Promise<Result<string, AurumError>> {
        try {
            const resp = await this.axios.get("/pk")
            const pk = resp.data as PublicKey;
            this.publicKey = pk.public_key;

            return ok(pk.public_key);
        } catch (error) {
            return err(error.response.data as AurumError);
        }
    }

    async Verify(token: string): Promise<Result<Claims, AurumError>> {
        if (this.publicKey === null) {
            const res = await this.GetPublicKey()
            if (res.isErr()) {
                return err(res.error)
            } else {
                this.publicKey = res.value
            }
        }

        const res = verifyJwt(token, this.publicKey)
        if (res.isErr()) {
            return err({
                Message: res.error,
                Code: ErrorCode.Unauthorized
            })
        }

        return ok(res.value)
    }


    // Login makes a request to log in a user. It returns either an error or null.
    // In the null case login was succesful and the Aurum client succesfully saved the token.
    async Login(user: User): Promise<Result<TokenPair, AurumError>> {
        try {
            const resp = await this.axios.post("/login", user)
            return ok(resp.data as TokenPair);
        } catch (error) {
            return err(error.response.data as AurumError);
        }
    }

    async Register(user: User): Promise<Result<null, AurumError>>{
        try {
            await this.axios.post("/signup", user)
            return ok(null);
        } catch (error) {
            return err(error.response.data as AurumError);
        }
    }

    async Refresh(tokenpair: TokenPair): Promise<Result<TokenPair, AurumError>> {
         try {
            const resp = await this.axios.post("/refresh", tokenpair)
            return ok(resp.data as TokenPair)
        } catch (error) {
            return err(error.response.data as AurumError)
        }
    }

    async GetUserInfo(tokenpair: TokenPair): Promise<Result<User, AurumError>> {
        const config: AxiosRequestConfig = {
            headers: {
                Authorization: "Bearer " + tokenpair.login_token
            }
        };

        try {
            const resp = await this.axios.get("/user", config)
            return ok(resp.data as User)
        } catch (error) {
            return err(error.response.data as AurumError)
        }
    }

    async UpdateUser(tokenpair: TokenPair, user: User): Promise<Result<User, AurumError>> {
        const config: AxiosRequestConfig = {
            headers: {
                Authorization: "Bearer " + tokenpair.login_token
            }
        };

        try {
            const resp = await this.axios.post("/user", user, config)
            return ok(resp.data as User)
        } catch (error) {
            return err(error.response.data as AurumError)
        }
    }

    async GetApplicationsForUser(tokenpair: TokenPair, user: User): Promise<Result<ApplicationWithRole[], AurumError>> {
        const config: AxiosRequestConfig = {
            headers: {
                Authorization: "Bearer " + tokenpair.login_token
            }
        };

        try {
            const resp = await this.axios.get("/application/" + user.username, config)
            return ok(resp.data as ApplicationWithRole[])
        } catch (error) {
            return err(error.response.data as AurumError)
        }
    }
}