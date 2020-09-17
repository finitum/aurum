use crate::error::AurumError;
use crate::token::TokenPair;
use crate::user::User;
use reqwest::blocking::Client;
use serde::{Deserialize, Serialize};
use std::ops::Range;
use url::Url;

/// Logs in the specified user using username and password
pub(crate) fn login(base_url: &Url, client: &Client, user: &User) -> Result<TokenPair, AurumError> {
    let resp = client.post(base_url.join("login")?).json(user).send()?;

    if resp.status().is_success() {
        Ok(resp.json()?)
    } else {
        Err(resp.status().into())
    }
}

/// Creates a new user using specified username and password
pub(crate) fn signup(base_url: &Url, client: &Client, user: &User) -> Result<(), AurumError> {
    let resp = client.post(base_url.join("signup")?).json(user).send()?;

    if resp.status().is_success() {
        Ok(())
    } else {
        Err(resp.status().into())
    }
}

#[derive(Serialize)]
pub(crate) struct RefreshRequest<'a> {
    pub(crate) refresh_token: &'a str,
}

#[derive(Serialize,Deserialize)]
pub(crate) struct RefreshResponse {
    pub(crate) login_token: String,
}

/// Refreshes the login token with the given refresh token
pub(crate) fn refresh<'a>(
    base_url: &Url,
    client: &Client,
    refresh_token: &RefreshRequest<'a>,
) -> Result<RefreshResponse, AurumError> {
    let resp = client
        .post(base_url.join("refresh")?)
        .json(refresh_token)
        .send()?;

    if resp.status().is_success() {
        Ok(resp.json()?)
    } else {
        Err(resp.status().into())
    }
}

pub(crate) fn refresh_tp(
    base_url: &Url,
    client: &Client,
    mut token_pair: TokenPair,
) -> Result<TokenPair, AurumError> {
    token_pair.login_token = String::new();

    let resp = client
        .post(base_url.join("refresh")?)
        .json(&token_pair)
        .send()?;

    if resp.status().is_success() {
        let new: TokenPair = resp.json()?;
        token_pair.login_token = new.login_token;

        Ok(token_pair)
    } else {
        Err(resp.status().into())
    }
}

#[derive(Serialize, Deserialize)]
pub(crate) struct PublicKeyResponse {
    pub(crate) public_key: String,
}

/// Gets the public key from the specified server
pub(crate) fn pk(base_url: &Url, client: &Client) -> Result<PublicKeyResponse, AurumError> {
    let resp = client.get(base_url.join("pk")?).send()?;

    if !resp.status().is_success() {
        Err(resp.status().into())
    } else {
        Ok(resp.json()?)
    }
}

// -- Authenticated Routes --

/// Gets the user struct of the current user
pub(crate) fn get_user(base_url: &Url, client: &Client, tokens: &TokenPair) -> Result<User, AurumError> {
    let bearer = format!("Bearer {}", tokens.login_token);

    let resp = client
        .get(base_url.join("user")?)
        .header("Authorization", bearer)
        .send()?;

    if !resp.status().is_success() {
        Err(resp.status().into())
    } else {
        Ok(resp.json()?)
    }
}

/// Update the user by providing a new user object, admins can change other users.
pub(crate) fn update_user(
    base_url: &Url,
    client: &Client,
    tokens: TokenPair,
    user: &User,
) -> Result<User, AurumError> {
    let bearer = format!("Bearer {}", tokens.login_token);

    let resp = client
        .put(base_url.join("me")?)
        .header("Authorization", bearer)
        .json(user)
        .send()?;

    if !resp.status().is_success() {
        Err(resp.status().into())
    } else {
        Ok(resp.json()?)
    }
}

// -- Admin Routes --

/// Gets a list of users limited by the specified range
pub(crate) fn users(
    base_url: &Url,
    client: &Client,
    tokens: TokenPair,
    range: Range<usize>,
) -> Result<Vec<User>, AurumError> {
    let bearer = format!("Bearer {}", tokens.login_token);

    let mut url = base_url.join("user")?;
    url.query_pairs_mut()
        .append_pair("start", range.start.to_string().as_str())
        .append_pair("end", range.end.to_string().as_str());

    let resp = client.get(url).header("Authorization", bearer).send()?;

    if !resp.status().is_success() {
        Err(resp.status().into())
    } else {
        Ok(resp.json()?)
    }
}
