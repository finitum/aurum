use crate::user::{User, AuthenticatedUser};
use crate::token::TokenPair;
use crate::error::{AurumError, Code};
use reqwest::blocking::Client;
use serde::{Serialize, Deserialize};

pub(crate) fn login(base_url: &str, client: &Client, user: &User) -> Result<TokenPair, AurumError> {
    let resp = client
        .post(&format!("{}/login", base_url))
        .json(user)
        .send()
        .map_err(|e| {
            AurumError::code(format!("A connection error occurred while logging in: {:?}", e), Code::ConnectionError)
        })?;

    if resp.status().is_success() {
        Ok(resp.json()?)
    } else {
        Err(resp.status().into())
    }
}

pub(crate) fn signup(base_url: &str, client: &Client, user: &User) -> Result<(), AurumError> {
    let resp = client
        .post(&format!("{}/signup", base_url))
        .json(user)
        .send()
        .map_err(|e| {
            AurumError::code(format!("A connection error occurred while signing up: {:?}", e), Code::ConnectionError)
        })?;

    if resp.status().is_success() {
        Ok(())
    } else {
        Err(resp.status().into())
    }
}

#[derive(Serialize)]
pub(crate) struct RefreshRequest<'a> {
    pub(crate) refresh_token: &'a str
}

#[derive(Deserialize)]
pub(crate) struct RefreshResponse {
    pub(crate) login_token: String
}

pub(crate) fn refresh<'a>(base_url: &str, client: &Client, refresh_token: &RefreshRequest<'a>) -> Result<RefreshResponse, AurumError> {
    let resp = client
        .post(&format!("{}/refresh", base_url))
        .json(refresh_token)
        .send()
        .map_err(|e| {
            AurumError::code(format!("A connection error occurred while refreshing the tokens: {:?}", e), Code::ConnectionError)
        })?;

    if resp.status().is_success() {
        Ok(resp.json()?)
    } else {
        Err(resp.status().into())
    }
}
