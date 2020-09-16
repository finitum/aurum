use serde::{Serialize, Deserialize};
use crate::token::TokenPair;
use crate::Aurum;
use crate::requests;
use crate::error::AurumError;
use crate::requests::{RefreshRequest, RefreshResponse};
use crate::error::Code;

#[repr(u8)]
#[derive(Serialize, Deserialize, Debug, PartialOrd, PartialEq)]
pub enum Role {
    User = 0,
    Admin = 1,
}

impl Default for Role {
    fn default() -> Self {
        Role::User
    }
}

#[derive(Serialize, Deserialize, Default, Debug, PartialEq)]
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
    token_pair: TokenPair
}

impl AuthenticatedUser {
    pub(crate) fn new(mut user: User, token_pair: TokenPair) -> Self {
        // CLEAR the password to make sure it can never be read again accidentally or leak
        user.password = String::new();

        Self {
            user,
            token_pair
        }
    }

    pub fn username(&self) -> &String {
        &self.user.username
    }

    pub fn email(&self) -> &String {
        &self.user.email
    }

    pub fn role(&self) -> &Role {
        &self.user.role
    }

    pub fn blocked(&self) -> bool {
        self.user.blocked
    }

    pub fn token(&self) -> &str {
        &self.token_pair.login_token
    }

    //---

    pub fn refresh_tokens(&mut self, aurum: &Aurum) -> Result<&str, AurumError> {
        let req = RefreshRequest {
            refresh_token: &self.token_pair.refresh_token,
        };

        let RefreshResponse {login_token} = requests::refresh(&aurum.base_url, &aurum.client, &req)?;

        self.token_pair.login_token = login_token;

        match self.token_pair.verify_tokens(&aurum.server_public_key) {
            true => Ok(&self.token_pair.login_token),
            _ => Err(AurumError::code("Received tokens not valid", Code::InvalidJWTToken))
        }
    }

    pub fn check_token(&mut self, aurum: &Aurum) -> Result<&str, AurumError> {
        if !self.token_pair.verify_tokens(&aurum.server_public_key) {
            self.refresh_tokens(aurum)
        } else {
            Ok(&self.token_pair.login_token)
        }
    }
}
