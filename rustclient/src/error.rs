use reqwest::StatusCode;
use der_parser::error::BerError;

#[repr(u8)]
#[derive(Debug, PartialOrd, PartialEq)]
pub enum Code {
    Unknown,
    InvalidCredentials,
    ServerError,
    InvalidJWTToken,
    ReqwestError,
    InvalidPEM,
    ConnectionError,
}

#[derive(Debug)]
pub struct AurumError {
    pub message: String,
    pub code: Code,
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

impl From<BerError> for AurumError {
    fn from(e: BerError) -> Self {
        AurumError::code(e.to_string(), Code::InvalidPEM)
    }
}

impl From<der_parser::nom::Err<der_parser::error::BerError>> for AurumError {
    fn from(e: der_parser::nom::Err<BerError>) -> Self {
        AurumError::code(e.to_string(), Code::InvalidPEM)
    }
}
