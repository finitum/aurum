use std::time::Duration;

use jwt_simple::algorithms::Ed25519PublicKey;
use reqwest::blocking::{Client, ClientBuilder};
use reqwest::header::HeaderMap;
use reqwest::Url;

use error::AurumError;
use user::User as InternalUser;

use crate::user::AuthenticatedUser;
pub use crate::user::AuthenticatedUser as User;
pub use crate::user::Role;

mod error;
mod requests;
mod token;
mod user;

#[derive(Debug)]
pub struct Aurum {
    base_url: Url,
    client: Client,

    // JWT Public key
    server_public_key: Ed25519PublicKey,
}

impl Aurum {
    pub fn new(base_url: String) -> Result<Self, AurumError> {
        log::info!("creating Aurum client at base url {}", base_url);
        let mut headers = HeaderMap::new();
        headers.insert(
            "Content-Type",
            "application/json".parse().map_err(AurumError::new)?,
        );

        let base_url = base_url
            .parse()
            .map_err(|e| AurumError::new(format!("failed to parse url: {:?}", e)))?;

        let client = ClientBuilder::new()
            .timeout(Duration::from_secs(5))
            .connect_timeout(Duration::from_secs(3))
            .default_headers(headers)
            .build()
            .map_err(|e| AurumError::new(format!("failed to create http client: {:?}", e)))?;

        let server_public_key = Self::get_pk(&client, &base_url)?;

        Ok(Self {
            base_url,
            client,
            server_public_key,
        })
    }

    /// Retrieves the server's public key and parses it into an [Ed25519PublicKey]
    fn get_pk(client: &Client, base_url: &Url) -> Result<Ed25519PublicKey, AurumError> {
        log::info!("requesting public key from Aurum server");

        let pk = requests::pk(base_url, client)?;

        token::pem_to_key(pk.public_key.as_str())
    }

    /// Logs in the user with Aurum.
    pub fn login(
        &mut self,
        username: String,
        password: String,
    ) -> Result<AuthenticatedUser, AurumError> {
        log::info!("logging in as {}", username);

        let user = InternalUser {
            username,
            password,
            ..Default::default()
        };

        let tp = requests::login(&self.base_url, &self.client, &user)?;

        Ok(AuthenticatedUser::new(user, tp))
    }

    /// Creates an account using the provided information and automatically logs in afterwards.
    pub fn signup(
        &mut self,
        username: String,
        email: String,
        password: String,
    ) -> Result<AuthenticatedUser, AurumError> {
        let user = InternalUser {
            username,
            password,
            email,
            ..Default::default()
        };
        requests::signup(&self.base_url, &self.client, &user)?;
        self.login(user.username.clone(), user.password)
    }
}

#[cfg(test)]
pub(crate) mod test_constants {
    pub const PUBLIC_TEST_KEY: &str = "-----BEGIN PUBLIC KEY-----\nMCowBQYDK2VwAyEAcYZfIh84oxMzA4bFmVFPNsSBCDn6D4nJiTTXsM46WGg=\n-----END PUBLIC KEY-----";
    pub const PUBLIC_TEST_KEY_B64: &str = "cYZfIh84oxMzA4bFmVFPNsSBCDn6D4nJiTTXsM46WGg=";
    pub const SECRET_TEST_KEY_B64: &str =
        "ovjfGUTfVkSQ6AP0qdFX7Z20FFHCPvDpKu5CeXXzVdRxhl8iHzijEzMDhsWZUU82xIEIOfoPicmJNNewzjpYaA==";
}

#[cfg(test)]
mod tests {
    use httpmock::{Mock, MockServer};
    use httpmock::Method;
    use jwt_simple::coarsetime::Duration;
    use jwt_simple::prelude::*;
    use reqwest::StatusCode;

    use crate::error::Code;
    use crate::token::{Claims as CustomClaims, generate_valid_tokenpair};
    use crate::token::TokenPair;
    use crate::user::Role;

    use super::*;
    use crate::test_constants::*;


    #[test]
    fn test_login() {
        let mock_server = MockServer::start();

        let (_, tp) = generate_valid_tokenpair("yeet");

        let pk = requests::PublicKeyResponse {
            public_key: PUBLIC_TEST_KEY.to_owned(),
        };

        let user = InternalUser {
            username: "user".to_string(),
            password: "pass".to_string(),
            ..Default::default()
        };

        let pk_mock = Mock::new()
            .expect_method(Method::GET)
            .expect_path("/pk")
            .return_status(200)
            .return_json_body(&pk)
            .create_on(&mock_server);

        let login_mock = Mock::new()
            .expect_method(Method::POST)
            .expect_path("/login")
            .expect_json_body(&user)
            .return_status(200)
            .return_json_body(&tp)
            .create_on(&mock_server);

        let url = format!("http://{}", mock_server.address());

        let mut au = Aurum::new(url).unwrap();
        let auth_user = au.login(user.username.clone(), user.password).unwrap();
        assert_eq!(auth_user.username(), &user.username);

        assert_eq!(
            au.server_public_key.to_bytes(),
            base64::decode(PUBLIC_TEST_KEY_B64).unwrap()
        );
        assert_eq!(auth_user.token(), tp.login_token);

        assert_eq!(pk_mock.times_called(), 1);
        assert_eq!(login_mock.times_called(), 1);
    }

    #[test]
    fn test_login_failure() {
        test_login_failure_helper(reqwest::StatusCode::UNAUTHORIZED, Code::Unauthorized);
        test_login_failure_helper(
            reqwest::StatusCode::INTERNAL_SERVER_ERROR,
            Code::ServerError,
        );
        test_login_failure_helper(reqwest::StatusCode::IM_A_TEAPOT, Code::Unknown);
    }

    fn test_login_failure_helper(code: reqwest::StatusCode, expected_error_code: Code) {
        let mock_server = MockServer::start();

        let (_, tp) = generate_valid_tokenpair("yeet");

        let pk = requests::PublicKeyResponse {
            public_key: PUBLIC_TEST_KEY.to_owned(),
        };

        let user = InternalUser {
            username: "user".to_string(),
            password: "pass".to_string(),
            ..Default::default()
        };

        let pk_mock = Mock::new()
            .expect_method(Method::GET)
            .expect_path("/pk")
            .return_status(200)
            .return_json_body(&pk)
            .create_on(&mock_server);

        let login_mock = Mock::new()
            .expect_method(Method::POST)
            .expect_path("/login")
            .expect_json_body(&user)
            .return_status(code.as_u16() as usize)
            .return_json_body(&tp)
            .create_on(&mock_server);

        let url = format!("http://{}", mock_server.address());

        let mut au = Aurum::new(url).unwrap();
        let err = au.login(user.username.clone(), user.password).unwrap_err();
        assert_eq!(err.code, expected_error_code);

        assert_eq!(pk_mock.times_called(), 1);
        assert_eq!(login_mock.times_called(), 1);
    }

    #[test]
    fn test_signup() {
        let mock_server = MockServer::start();

        let (_, tp) = generate_valid_tokenpair("yeet");

        let pk = requests::PublicKeyResponse {
            public_key: PUBLIC_TEST_KEY.to_owned(),
        };

        let user = InternalUser {
            username: "user".to_string(),
            password: "pass".to_string(),
            email: "email@example.com".to_string(),
            ..Default::default()
        };

        let pk_mock = Mock::new()
            .expect_method(Method::GET)
            .expect_path("/pk")
            .return_status(200)
            .return_json_body(&pk)
            .create_on(&mock_server);

        let signup_mock = Mock::new()
            .expect_method(Method::POST)
            .expect_path("/signup")
            .expect_json_body(&user)
            .return_status(StatusCode::CREATED.as_u16() as usize)
            .create_on(&mock_server);

        let login_mock = Mock::new()
            .expect_method(Method::POST)
            .expect_path("/login")
            .expect_json_body(&InternalUser {
                username: user.username.clone(),
                password: user.password.clone(),
                ..Default::default()
            })
            .return_status(200)
            .return_json_body(&tp)
            .create_on(&mock_server);

        let url = format!("http://{}", mock_server.address());

        let mut au = Aurum::new(url).unwrap();
        let auth_user = au.signup(user.username.clone(), user.email.clone(), user.password).unwrap();
        assert_eq!(auth_user.username(), &user.username);
        assert_eq!(auth_user.token(), tp.login_token);

        assert_eq!(
            au.server_public_key.to_bytes(),
            base64::decode(PUBLIC_TEST_KEY_B64).unwrap()
        );

        assert_eq!(pk_mock.times_called(), 1);
        assert_eq!(login_mock.times_called(), 1);
        assert_eq!(signup_mock.times_called(), 1);
    }

    #[test]
    fn test_signup_failure() {
        test_signup_failure_helper(reqwest::StatusCode::CONFLICT, Code::Conflict);
        test_signup_failure_helper(reqwest::StatusCode::UNPROCESSABLE_ENTITY, Code::InsufficientPassword);
        test_signup_failure_helper(reqwest::StatusCode::INTERNAL_SERVER_ERROR, Code::ServerError);
        test_signup_failure_helper(reqwest::StatusCode::IM_A_TEAPOT, Code::Unknown);
    }

    fn test_signup_failure_helper(code: reqwest::StatusCode, expected_error_code: Code) {
        let mock_server = MockServer::start();

        let (_, tp) = generate_valid_tokenpair("yeet");

        let pk = requests::PublicKeyResponse {
            public_key: PUBLIC_TEST_KEY.to_owned(),
        };

        let user = InternalUser {
            username: "user".to_string(),
            password: "pass".to_string(),
            email: "email@example.com".to_string(),
            ..Default::default()
        };

        let pk_mock = Mock::new()
            .expect_method(Method::GET)
            .expect_path("/pk")
            .return_status(200)
            .return_json_body(&pk)
            .create_on(&mock_server);

        let signup_mock = Mock::new()
            .expect_method(Method::POST)
            .expect_path("/signup")
            .expect_json_body(&user)
            .return_status(code.as_u16() as usize)
            .create_on(&mock_server);

        let url = format!("http://{}", mock_server.address());

        let mut au = Aurum::new(url).unwrap();
        let err = au.signup(user.username, user.email, user.password).unwrap_err();
        assert_eq!(err.code, expected_error_code);

        assert_eq!(pk_mock.times_called(), 1);
        assert_eq!(signup_mock.times_called(), 1);
    }
}
