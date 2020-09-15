mod user;
mod error;
mod token;

use reqwest::blocking::{ClientBuilder, Client};
use std::time::Duration;
use user::User;
use error::AurumError;
use reqwest::header::HeaderMap;
use crate::token::TokenPair;
use serde_json::Value;

struct Aurum {
    base_url: String,
    client: Client,
    decoding_key: jsonwebtoken::DecodingKey<'static>
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

        let pk = Self::get_pk(&base_url, &client)?;

        let decoding_key = jsonwebtoken::DecodingKey::from_ec_pem(pk.as_ref())
            .map_err(|e| AurumError::new(format!("failed to parse jwt pem: {:?}", e)))?
            .into_static();

        Ok(Self {
            base_url,
            client,
            decoding_key,
        })
    }

    fn get_pk(base_url: &str, client: &Client) -> Result<String, AurumError>  {
        let res = client.get(&format!("{}/pk", base_url))
            .send()
            .map_err(|_| AurumError::new("requesting aurum public key failed"))?;

        if !res.status().is_success() {
            return Err(res.status().into())
        }

        Ok(if let Value::Object(obj) = res.json()
            .map_err(|e| AurumError::new(format!("failed to decode json while requesting aurum public key: {:?}", e)))?
        {
            obj.get("public_key")
                .ok_or_else(|| AurumError::new("failed to decode json while requesting aurum public key"))?
                .to_string()
        } else {
            return Err(AurumError::new("failed to decode json while requesting aurum public key"));
        })
    }


    pub fn login(&self, username: String, password: String) -> Result<TokenPair, AurumError> {
        let user = User{username, password, ..Default::default()};
        
        let resp = self.client
            .post(&format!("{}/login", &self.base_url))
            .json(&user)
            .send()
            .map_err(|e| {
                AurumError::new(format!("An error occurred while logging in: {:?}", e))
            })?;

        if resp.status().is_success() {
            let tp = resp.json()?;
            Ok(tp)
        } else {
            Err(resp.status().into())
        }
    }
}
