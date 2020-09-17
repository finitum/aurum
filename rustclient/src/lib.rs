mod error;
mod requests;
mod token;
mod user;

use crate::user::AuthenticatedUser;
use error::AurumError;
use jwt_simple::algorithms::Ed25519PublicKey;
use reqwest::blocking::{Client, ClientBuilder};
use reqwest::header::HeaderMap;
use reqwest::Url;
use serde::{Deserialize, Serialize};
use std::time::Duration;
use user::User as InternalUser;

pub use crate::user::AuthenticatedUser as User;
pub use crate::user::Role;

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
    /// TODO: move to requests
    fn get_pk(client: &Client, base_url: &Url) -> Result<Ed25519PublicKey, AurumError> {
        log::info!("requesting public key from Aurum server");

        let pk = requests::pk(base_url, client)?;

        token::pem_to_key(pk.public_key.as_str())
    }

    /// Logs in the user with Aurum. TODO: more comments
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

    /// TODO: docs
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

    // pub fn user(&mut self, username: String) {
    //     requests::user(&self.base_url, &self.client, &user)?;
    //
    // }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::error::Code;
    use crate::token::Claims as CustomClaims;
    use crate::token::TokenPair;
    use crate::user::Role;
    use httpmock::Method;
    use httpmock::{Mock, MockServer};
    use jwt_simple::coarsetime::Duration;
    use jwt_simple::prelude::*;

    const PUBLIC_TEST_KEY: &str = "-----BEGIN PUBLIC KEY-----\nMCowBQYDK2VwAyEAcYZfIh84oxMzA4bFmVFPNsSBCDn6D4nJiTTXsM46WGg=\n-----END PUBLIC KEY-----";
    const PUBLIC_TEST_KEY_B64: &str = "cYZfIh84oxMzA4bFmVFPNsSBCDn6D4nJiTTXsM46WGg=";
    const SECRET_TEST_KEY_B64: &str =
        "ovjfGUTfVkSQ6AP0qdFX7Z20FFHCPvDpKu5CeXXzVdRxhl8iHzijEzMDhsWZUU82xIEIOfoPicmJNNewzjpYaA==";

    fn generate_valid_tokenpair(username: &str) -> (Ed25519KeyPair, TokenPair) {
        let key = Ed25519KeyPair::from_bytes(base64::decode(SECRET_TEST_KEY_B64).unwrap().as_ref())
            .unwrap();
        let lc = CustomClaims::new(username.to_owned(), Role::default(), false);
        let rc = CustomClaims::new(username.to_owned(), Role::default(), true);
        let lclaims = Claims::with_custom_claims(lc, Duration::from_hours(2));
        let rclaims = Claims::with_custom_claims(rc, Duration::from_hours(2));
        let login_token = key.sign(lclaims).unwrap();
        let refresh_token = key.sign(rclaims).unwrap();

        // Sanity check verify
        assert!(key
            .public_key()
            .verify_token::<CustomClaims>(&login_token, None)
            .is_ok());
        assert!(key
            .public_key()
            .verify_token::<CustomClaims>(&refresh_token, None)
            .is_ok());

        (key, TokenPair::from_tokens(login_token, refresh_token))
    }

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
        test_login_failure_helper(reqwest::StatusCode::UNAUTHORIZED, Code::InvalidCredentials);
        test_login_failure_helper(
            reqwest::StatusCode::INTERNAL_SERVER_ERROR,
            Code::ServerError,
        );
        test_login_failure_helper(reqwest::StatusCode::from_u16(420).unwrap(), Code::Unknown);
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
}
