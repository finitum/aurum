use crate::user::User;
use crate::token::TokenPair;
use crate::error::{AurumError, Code};
use reqwest::blocking::Client;
use serde::{Serialize, Deserialize};
use reqwest::Url;

pub(crate) fn join(url: &Url, path: &str) -> Result<Url, AurumError>{
    url.join(path).map_err(|e| AurumError::new(format!("failed to parse url: {:?}", e)))
}

pub(crate) fn login(base_url: &Url, client: &Client, user: &User) -> Result<TokenPair, AurumError> {
    let resp = client
        .post(join(base_url, "login")?)
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

pub(crate) fn signup(base_url: &Url, client: &Client, user: &User) -> Result<(), AurumError> {
    let resp = client
        .post(join(base_url, "signup")?)
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

pub(crate) fn get_user(base_url: &Url, client: &Client, user: &User) -> Result<(), AurumError> {
    let resp = client
        .post(join(base_url, "signup")?)
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

pub(crate) fn refresh<'a>(base_url: &Url, client: &Client, refresh_token: &RefreshRequest<'a>) -> Result<RefreshResponse, AurumError> {
    let resp = client
        .post(join(base_url, "refresh")?)
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
