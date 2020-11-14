export interface User {
    username: string
    password: string
    email: string
}

export interface Application {
    name: string
}


export enum ErrorCode {
    ServerError = 1,
    InvalidRequest,
    WeakPassword,
    Unauthorized,
}

export interface AurumError {
    Message: string
    Code:    ErrorCode
}

export interface TokenPair {
    login_token: string,
    refresh_token: string | null,
}
