mod user;
mod error;
mod token;
mod requests;

use reqwest::blocking::{ClientBuilder, Client};
use std::time::Duration;
use user::User;
use error::AurumError;
use reqwest::header::HeaderMap;
use jwt_simple::algorithms::Ed25519PublicKey;
use serde::{Serialize, Deserialize};
use crate::error::Code;
use crate::user::AuthenticatedUser;

#[derive(Debug)]
pub struct Aurum {
    base_url: String,
    client: Client,

    // JWT Public key
    server_public_key: Ed25519PublicKey,
}

#[derive(Serialize,Deserialize)]
struct PublicKey {
    public_key: String
}

impl Aurum {
    pub fn new(base_url: String) -> Result<Self, AurumError> {
        let mut headers = HeaderMap::new();
        headers.insert("Content-Type", "application/json".parse().map_err(AurumError::new)?);

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
    fn get_pk(client: &Client, base_url: &str) -> Result<Ed25519PublicKey, AurumError>  {
        let res = client.get(&format!("{}/pk", base_url))
            .send()
            .map_err(|_| AurumError::new("requesting aurum public key failed"))?;

        if !res.status().is_success() {
            return Err(res.status().into())
        }

        let json: PublicKey = res.json()
            .map_err(|e| AurumError::new(format!("failed to decode json while requesting aurum public key: {:?}", e)))?;

        token::pem_to_key(json.public_key.as_str())
    }

    
    /// Logs in the user with Aurum. TODO: more comments
    pub fn login(&mut self, username: String, password: String) -> Result<AuthenticatedUser, AurumError> {
        let user = User {
            username,
            password,
            ..Default::default()
        };

        let tp = requests::login(&self.base_url, &self.client, &user)?;

        Ok(AuthenticatedUser::new(user, tp))
    }

    pub fn signup(&mut self, username: String, password: String, email: String) -> Result<AuthenticatedUser, AurumError> {
        let user = User{username, password, email, ..Default::default()};
        requests::signup(&self.base_url, &self.client, &user)?;
        self.login(user.username.clone(), user.password)
    }
}

#[cfg(test)]
mod tests {
    use httpmock::{MockServer, Mock};
    use httpmock::Method;
    use crate::token::{TokenPair};
    use crate::user::Role;
    use crate::token::Claims as CustomClaims;
    use jwt_simple::prelude::*;
    use super::*;
    use jwt_simple::coarsetime::Duration;

    const PUBLIC_TEST_KEY: &str = "-----BEGIN PUBLIC KEY-----\nMCowBQYDK2VwAyEAcYZfIh84oxMzA4bFmVFPNsSBCDn6D4nJiTTXsM46WGg=\n-----END PUBLIC KEY-----";
    const PUBLIC_TEST_KEY_B64: &str = "cYZfIh84oxMzA4bFmVFPNsSBCDn6D4nJiTTXsM46WGg=";
    const SECRET_TEST_KEY_B64: &str = "ovjfGUTfVkSQ6AP0qdFX7Z20FFHCPvDpKu5CeXXzVdRxhl8iHzijEzMDhsWZUU82xIEIOfoPicmJNNewzjpYaA==";

    fn generate_valid_tokenpair(username: &str) -> (Ed25519KeyPair, TokenPair) {
        let key = Ed25519KeyPair::from_bytes(base64::decode(SECRET_TEST_KEY_B64).unwrap().as_ref()).unwrap();
        let lc = CustomClaims::new(username.to_owned(), Role::default(), false);
        let rc = CustomClaims::new(username.to_owned(), Role::default(), true);
        let lclaims = Claims::with_custom_claims(lc, Duration::from_hours(2));
        let rclaims = Claims::with_custom_claims(rc, Duration::from_hours(2));
        let login_token = key.sign(lclaims).unwrap();
        let refresh_token = key.sign(rclaims).unwrap();

        // Sanity check verify
        assert!(key.public_key().verify_token::<CustomClaims>(&login_token, None).is_ok());
        assert!(key.public_key().verify_token::<CustomClaims>(&refresh_token, None).is_ok());

        (key, TokenPair::from_tokens(login_token, refresh_token))
    }

    #[test]
    fn test_login() {
        let mock_server = MockServer::start();

        let (key, tp) = generate_valid_tokenpair("yeet");

        let pk = PublicKey{
            public_key: PUBLIC_TEST_KEY.to_owned()
        };

        let user = User{
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


        assert_eq!(au.server_public_key.to_bytes(), base64::decode(PUBLIC_TEST_KEY_B64).unwrap());
        assert_eq!(auth_user.token(), tp.login_token);

        assert_eq!(pk_mock.times_called(), 1);
        assert_eq!(login_mock.times_called(), 1);
    }

    #[test]
    fn test_login_failure() {
        test_login_failure_helper(reqwest::StatusCode::UNAUTHORIZED, Code::InvalidCredentials);
        test_login_failure_helper(reqwest::StatusCode::INTERNAL_SERVER_ERROR, Code::ServerError);
        test_login_failure_helper(reqwest::StatusCode::from_u16(420).unwrap(), Code::Unknown);
    }

    fn test_login_failure_helper(code: reqwest::StatusCode, expected_error_code: Code) {
        let mock_server = MockServer::start();

        let (key, tp) = generate_valid_tokenpair("yeet");

        let pk = PublicKey{
            public_key: PUBLIC_TEST_KEY.to_owned()
        };

        let user = User{
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
