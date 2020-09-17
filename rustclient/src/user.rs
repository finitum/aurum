use crate::error::AurumError;
use crate::error::Code;
use crate::requests;
use crate::requests::{RefreshRequest, RefreshResponse};
use crate::token::TokenPair;
use crate::Aurum;
use serde::{Deserialize, Serialize};
use serde_repr::*;

#[repr(u8)]
#[derive(Serialize_repr, Deserialize_repr, Debug, PartialOrd, PartialEq, Clone)]
pub enum Role {
    User = 0,
    Admin = 1,
}

impl Default for Role {
    fn default() -> Self {
        Role::User
    }
}

#[derive(Serialize, Deserialize, Default, Debug, PartialEq, Clone)]
pub(crate) struct User {
    pub(crate) username: String,
    pub(crate) password: String,
    pub(crate) email: String,
    pub(crate) role: Role,
    pub(crate) blocked: bool,
}

#[derive(Debug)]
pub struct AuthenticatedUser {
    user: User,
    token_pair: TokenPair,
}

impl AuthenticatedUser {
    pub(crate) fn new(mut user: User, token_pair: TokenPair) -> Self {
        // CLEAR the password to make sure it can never be read again accidentally or leak
        user.password = String::new();

        Self { user, token_pair }
    }

    /// Returns this users username
    pub fn username(&self) -> &String {
        &self.user.username
    }

    /// Returns this users email
    pub fn email(&self) -> &String {
        &self.user.email
    }

    /// Returns the [Role] of this user
    pub fn role(&self) -> &Role {
        &self.user.role
    }

    /// Returns true if this user is blocked and false otherwise
    pub fn blocked(&self) -> bool {
        self.user.blocked
    }

    /// Returns the login token for this user
    pub fn token(&self) -> &str {
        &self.token_pair.login_token
    }

    //---

    /// Refreshes the login token unconditionally
    /// See [check_token] for checking and refreshing
    pub(crate) fn refresh_tokens(&mut self, aurum: &Aurum) -> Result<&str, AurumError> {
        let req = RefreshRequest {
            refresh_token: &self.token_pair.refresh_token,
        };

        let RefreshResponse { login_token } =
            requests::refresh(&aurum.base_url, &aurum.client, &req)?;

        self.token_pair.login_token = login_token;

        match self.token_pair.verify_tokens(&aurum.server_public_key) {
            true => Ok(&self.token_pair.login_token),
            _ => Err(AurumError::code(
                "Received tokens not valid",
                Code::InvalidJWTToken,
            )),
        }
    }

    /// Checks if the login token needs to be refreshed and refreshes it if necessary.
    pub fn check_token(&mut self, aurum: &Aurum) -> Result<&str, AurumError> {
        if !self.token_pair.verify_tokens(&aurum.server_public_key) {
            self.refresh_tokens(aurum)
        } else {
            Ok(&self.token_pair.login_token)
        }
    }

    /// Grabs the latest user info and saves this into self.
    pub fn get_user(&mut self, aurum: &Aurum) -> Result<(), AurumError> {
        let updated_user = requests::get_user(&aurum.base_url, &aurum.client, &self.token_pair)?;

        if self.user.username == updated_user.username {
            self.user = updated_user;
            Ok(())
        } else {
            Err(AurumError::code("Received user does not match expectations", Code::InvalidResponse))
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::token::generate_valid_tokenpair;
    use httpmock::{Mock, Method, MockServer};
    use crate::test_constants::{PUBLIC_TEST_KEY, PUBLIC_TEST_KEY_B64};

    #[test]
    fn test_getters() {
        let u = User {
            username: "user".to_string(),
            password: "pass".to_string(),
            email: "email".to_string(),
            role: Role::Admin,
            blocked: true
        };

        let (_, tp) = generate_valid_tokenpair(u.username.as_str());

        let au = AuthenticatedUser::new(u.clone(), tp.clone());

        assert_eq!(au.username(), &u.username);
        assert_eq!(au.email(), &u.email);
        assert_eq!(au.role(), &u.role);
        assert_eq!(au.blocked(), u.blocked);
        assert_eq!(au.user.password, String::new());

        assert_eq!(au.token(), tp.login_token)
    }

    #[test]
    fn test_refresh() {
        let mock_server = MockServer::start();

        let pk = requests::PublicKeyResponse {
            public_key: PUBLIC_TEST_KEY.to_owned(),
        };

        let (_, newtp) = generate_valid_tokenpair("yeet");

        let resp = RefreshResponse{
            login_token: newtp.login_token.clone()
        };

        let pk_mock = Mock::new()
            .expect_method(Method::GET)
            .expect_path("/pk")
            .return_status(200)
            .return_json_body(&pk)
            .create_on(&mock_server);

        let refresh_mock = Mock::new()
            .expect_method(Method::POST)
            .expect_path("/refresh")
            .return_status(200)
            .return_json_body(&resp)
            .create_on(&mock_server);

        let url = format!("http://{}", mock_server.address());


        let (_, tp) = generate_valid_tokenpair("yeet");

        let au = Aurum::new(url).unwrap();
        let mut user = AuthenticatedUser{
            user: User{
                username: "yeet".to_string(),
                ..Default::default()
            },
            token_pair: tp
        };

        let res = user.refresh_tokens(&au).unwrap();
        assert_eq!(res, &newtp.login_token);
        assert_eq!(user.token_pair.login_token, newtp.login_token);

        assert_eq!(pk_mock.times_called(), 1);
        assert_eq!(refresh_mock.times_called(), 1);
    }

    #[test]
    fn test_get_user() {
        let mock_server = MockServer::start();

        let pk = requests::PublicKeyResponse {
            public_key: PUBLIC_TEST_KEY.to_owned(),
        };

        let orig = User{
            username: "user".to_string(),
            password: "pass".to_string(),
            email: "old-email".to_string(),
            ..Default::default()
        };

        let mut new = orig.clone();
        new.email = "new mail".to_string();

        let (_, tp) = generate_valid_tokenpair("yeet");

        let pk_mock = Mock::new()
            .expect_method(Method::GET)
            .expect_path("/pk")
            .return_status(200)
            .return_json_body(&pk)
            .create_on(&mock_server);

        let get_user_mock = Mock::new()
            .expect_method(Method::GET)
            .expect_path("/user")
            .expect_header("Authorization", &format!("Bearer {}", &tp.login_token))
            .return_status(200)
            .return_json_body(&new)
            .create_on(&mock_server);

        let url = format!("http://{}", mock_server.address());

        let au = Aurum::new(url).unwrap();
        let mut auth_user = AuthenticatedUser{
            user: orig,
            token_pair: tp
        };

        auth_user.get_user(&au).unwrap();
        assert_eq!(auth_user.user, new);

        assert_eq!(pk_mock.times_called(), 1);
        assert_eq!(get_user_mock.times_called(), 1);
    }

    #[test]
    fn test_get_user_failure() {
        let mock_server = MockServer::start();

        let pk = requests::PublicKeyResponse {
            public_key: PUBLIC_TEST_KEY.to_owned(),
        };

        let orig = User{
            username: "user".to_string(),
            password: "pass".to_string(),
            email: "old-email".to_string(),
            ..Default::default()
        };

        let mut new = orig.clone();
        new.username = "otheruser".to_string();

        let (_, tp) = generate_valid_tokenpair("yeet");

        let pk_mock = Mock::new()
            .expect_method(Method::GET)
            .expect_path("/pk")
            .return_status(200)
            .return_json_body(&pk)
            .create_on(&mock_server);

        let get_user_mock = Mock::new()
            .expect_method(Method::GET)
            .expect_path("/user")
            .expect_header("Authorization", &format!("Bearer {}", &tp.login_token))
            .return_status(200)
            .return_json_body(&new)
            .create_on(&mock_server);

        let url = format!("http://{}", mock_server.address());

        let au = Aurum::new(url).unwrap();
        let mut auth_user = AuthenticatedUser{
            user: orig.clone(),
            token_pair: tp
        };

        let err = auth_user.get_user(&au).unwrap_err();
        assert_eq!(err.code, Code::InvalidResponse);
        assert_eq!(auth_user.user, orig);


        assert_eq!(pk_mock.times_called(), 1);
        assert_eq!(get_user_mock.times_called(), 1);
    }
}
