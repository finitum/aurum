
export interface User {
    username: string
    password: string
    email: string
}

export interface Application {
    name: string
    allow_registration: string
}

export interface ApplicationWithRole extends Application {
    role: Role
}

export enum Role {
    User = 1,
    Admin
}

export enum ErrorCode {
    ServerError ,
    InvalidRequest,
    Duplicate,
    WeakPassword,
    Unauthorized,
}

export interface AurumError {
    Message: string
    Code:    ErrorCode
}

export interface TokenPair {
    login_token: string,
    refresh_token: string,
}
