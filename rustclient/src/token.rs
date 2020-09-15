use jsonwebtoken::DecodingKey;
use super::user;
use serde::{Serialize, Deserialize};
use crate::error::AurumError;
use crate::error::Code;

#[derive(Serialize, Deserialize)]
pub(crate) struct Claims {
    username: String,
    role: user::Role,

    refresh: bool,
    // Issued at
    iat: i64,
    // Expiry
    exp: i64,
    // Not Before
    nbf: i64,
}

#[derive(Serialize, Deserialize)]
pub(crate) struct TokenPair {
    #[serde(skip_serializing_if = "Option::is_none")]
    refresh_token: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    login_token: Option<String>,

    #[serde(skip)]
    claims: Option<Claims>,
}

impl TokenPair {
    pub fn claims<'a>(&'a mut self, dec: &DecodingKey) -> Result<&'a Claims, AurumError> {
        let token = self.login_token.as_ref()
            .ok_or_else(|| AurumError::code("no login token found in tokenpair", Code::InvalidJWTToken))?;


        if let Some(ref c) = self.claims {
            return Ok(c)
        }

        // decode
        let tokendata = jsonwebtoken::decode(&token, dec, &jsonwebtoken::Validation{
            leeway: 5,
            validate_exp: true,
            validate_nbf: true,
            ..Default::default()
        }).map_err(|e| AurumError::code(e, Code::InvalidJWTToken))?;

        self.claims = Some(tokendata.claims);

        // unwrap is safe, just look one line up
        Ok(self.claims.as_ref().unwrap())
    }
}
