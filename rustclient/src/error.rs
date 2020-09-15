use reqwest::StatusCode;

#[repr(u8)]
pub enum Code {
    Unknown,
    InvalidCredentials,
    ServerError,
    InvalidJWTToken,
    ReqwestError,
}

pub struct AurumError {
    message: String,
    code: Code,
}

impl AurumError {
    pub fn new<T: ToString>(s: T) -> Self {
        AurumError {
            message: s.to_string(),
            code: Code::Unknown
        }
    }

    pub fn code<T: ToString>(s: T, code: Code) -> Self {
        AurumError {
            message: s.to_string(),
            code,
        }
    }
}

impl From<StatusCode> for AurumError {
    fn from(status: StatusCode) -> Self {
        match status {
            StatusCode::UNAUTHORIZED => AurumError {
                message: status.to_string(),
                code: Code::InvalidCredentials
            },
            StatusCode::INTERNAL_SERVER_ERROR => AurumError {
                message: status.to_string(),
                code: Code::ServerError
            },
            _ => AurumError {
                message: status.to_string(),
                code: Code::Unknown
            }
        }
    }
}

impl From<reqwest::Error> for AurumError {
    fn from(e: reqwest::Error) -> Self {
        AurumError::code(e.to_string(),Code::ReqwestError)
    }
}
