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
    pub fn refresh_tokens(&mut self, aurum: &Aurum) -> Result<&str, AurumError> {
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
    /// See [refresh_token] for unconditional refresh
    pub fn check_token(&mut self, aurum: &Aurum) -> Result<&str, AurumError> {
        if !self.token_pair.verify_tokens(&aurum.server_public_key) {
            self.refresh_tokens(aurum)
        } else {
            Ok(&self.token_pair.login_token)
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::token::generate_valid_tokenpair;

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
}
