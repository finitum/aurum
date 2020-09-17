use super::user::Role;
use crate::error::AurumError;
use crate::error::Code;
use crate::error::Code::InvalidPEM;
use jwt_simple::prelude::*;
use serde::{Deserialize, Deserializer, Serialize};

#[derive(Serialize, Deserialize, Debug, PartialEq)]
pub(crate) struct Claims {
    username: String,
    role: Role,

    refresh: bool,
}

impl Claims {
    #[cfg(test)]
    pub(crate) fn new(username: String, role: Role, refresh: bool) -> Self {
        Self {
            username,
            role,
            refresh,
        }
    }
}

#[derive(Serialize, Deserialize, Default, Debug, Clone)]
pub(crate) struct TokenPair {
    pub(crate) login_token: String,
    #[serde(deserialize_with = "none_is_empty_string")]
    pub(crate) refresh_token: String,
}

fn none_is_empty_string<'de, D: Deserializer<'de>>(d: D) -> Result<String, D::Error> {
    Option::deserialize(d).map(Option::unwrap_or_default)
}

impl TokenPair {
    // pub(crate) fn from_tokens(login_token: String, refresh_token: String) -> Self {
    //     TokenPair {
    //         refresh_token,
    //         login_token,
    //     }
    // }

    /// Verifies the signatures on the two tokens inside. Sets the two claims fields. Returns false
    /// if the verification failed, but if true it makes sure the two claims fields are *NOT* None.
    ///
    /// Therefore, if later any of the claims fields are None for whatever reason, this is a key verification error.
    pub(crate) fn verify_tokens(&self, key: &Ed25519PublicKey) -> bool {
        let login_claims = key.verify_token::<Claims>(&self.login_token, None);
        let refresh_claims = key.verify_token::<Claims>(&self.login_token, None);

        matches!((login_claims, refresh_claims), (Ok(_), Ok(_)))
    }
}

const ED25519_OID: &str = "1.3.101.112";

pub(crate) fn pem_to_key(pem: &str) -> Result<Ed25519PublicKey, AurumError> {
    // Filter headers
    let b64 = pem
        .lines()
        .filter(|&i| !i.contains("BEGIN PUBLIC KEY") && !i.contains("END PUBLIC KEY"))
        .collect::<String>();

    // Base64 -> binary
    let binary = base64::decode(b64).map_err(|e| AurumError::code(e, InvalidPEM))?;

    // oid is nested 3 deep
    let (_, der) = der_parser::der::parse_der_recursive(&binary, 3)?;

    // Get content from der
    let content = der
        .as_sequence()
        .map_err(|e| AurumError::code(e, InvalidPEM))?;

    // Get and check the oid
    let oid = content
        .get(0)
        .ok_or_else(|| AurumError::code("Content Vec is empty", Code::InvalidPEM))?
        .as_sequence()?
        .get(0)
        .ok_or_else(|| AurumError::code("Oid Sequence is empty", Code::InvalidPEM))?
        .as_oid()?;

    if oid.to_id_string() != ED25519_OID {
        return Err(AurumError::code(
            "Unexpected Crypto Algorithm",
            Code::InvalidPEM,
        ));
    }

    // Retrieve key bytes
    let key = content
        .get(1)
        .ok_or_else(|| AurumError::code("Content Vec is empty", Code::InvalidPEM))?
        .as_bitstring()?
        .data;

    Ed25519PublicKey::from_bytes(key).map_err(|e| AurumError::code(e.to_string(), Code::InvalidPEM))
}

#[cfg(test)]
pub(crate) fn generate_valid_tokenpair(username: &str) -> (Ed25519KeyPair, TokenPair) {
    use jwt_simple::claims::Claims as jwtClaims;
    use crate::test_constants::SECRET_TEST_KEY_B64;
    use crate::Role;

    let key = Ed25519KeyPair::from_bytes(base64::decode(SECRET_TEST_KEY_B64).unwrap().as_ref())
        .unwrap();
    let lc = Claims::new(username.to_owned(), Role::default(), false);
    let rc = Claims::new(username.to_owned(), Role::default(), true);
    let lclaims = jwtClaims::with_custom_claims(lc, Duration::from_hours(2));
    let rclaims = jwtClaims::with_custom_claims(rc, Duration::from_hours(2));
    let login_token = key.sign(lclaims).unwrap();
    let refresh_token = key.sign(rclaims).unwrap();

    // Sanity check verify
    assert!(key
        .public_key()
        .verify_token::<Claims>(&login_token, None)
        .is_ok());
    assert!(key
        .public_key()
        .verify_token::<Claims>(&refresh_token, None)
        .is_ok());

    (key, TokenPair::from_tokens(login_token, refresh_token))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_der_parser_vs_reference_impl() {
        let pem = "-----BEGIN PUBLIC KEY-----\nMCowBQYDK2VwAyEAke7D81dEbGP9xiHsQ0/qIn/BYiphsY8qk3iSVBHTYXs=\n-----END PUBLIC KEY-----\n";

        let b64 = pem
            .lines()
            .filter(|&i| !i.contains("BEGIN PUBLIC KEY") && !i.contains("END PUBLIC KEY"))
            .collect::<String>();

        let binary = base64::decode(b64).unwrap();

        // key + metadata (oid)
        assert_eq!(binary.len(), 32 + 12);

        let (_, der) = der_parser::parse_der(&binary).unwrap();

        let content = der.content.as_sequence().unwrap();
        let oid = content
            .get(0)
            .unwrap()
            .as_sequence()
            .unwrap()
            .get(0)
            .unwrap()
            .as_oid()
            .unwrap();
        assert_eq!(oid.to_id_string(), "1.3.101.112");
        assert_eq!("1.3.101.112", ED25519_OID);

        let key = content.get(1).unwrap().as_bitstring().unwrap().data;

        assert_eq!(key.len(), 32);

        let key2 = pem_to_key(pem).unwrap().to_bytes();

        assert_eq!(key, key2.as_slice());
    }

    #[test]
    fn test_pem_to_key() {
        let pem = "-----BEGIN PUBLIC KEY-----\nMCowBQYDK2VwAyEAke7D81dEbGP9xiHsQ0/qIn/BYiphsY8qk3iSVBHTYXs=\n-----END PUBLIC KEY-----\n";
        let openssl_pem = "-----BEGIN PUBLIC KEY-----\nMCowBQYDK2VwAyEA1vpaMjzK5je2AWxuSvxWL7dkXC55HA7Wx/laIxdOb5M=\n-----END PUBLIC KEY-----\n";
        let go_pem = "-----BEGIN PUBLIC KEY-----\nMCowBQYDK2VwAyEAD0pl7Ix4FL+4CueSdaKQ2yP72CDgJXZd/XCoBF41U4A=\n-----END PUBLIC KEY-----\n";
        pem_to_key(pem).unwrap();
        pem_to_key(openssl_pem).unwrap();
        pem_to_key(go_pem).unwrap();
    }

    #[test]
    fn test_pem_to_key_matches() {
        let go_pem = "-----BEGIN PUBLIC KEY-----\nMCowBQYDK2VwAyEAuDduaauAxMocf7HCvWAUzgrwM83iYB/loufqmHsyCuI=\n-----END PUBLIC KEY-----";
        let go_b64 = "uDduaauAxMocf7HCvWAUzgrwM83iYB/loufqmHsyCuI=";
        let key = pem_to_key(go_pem).unwrap();
        let expected = base64::decode(go_b64).unwrap();

        assert_eq!(expected, key.to_bytes());
    }

    #[test]
    fn test_pem_rsa() {
        let rsa_pem = "-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAp6O3wt34sJc5yu5vo/3V
QDoHBTQ6QlYNUcblVPj6+4naulbrZC+NjDmJYAZsO4IsWQjG9JeUEzLN/9mcPw6s
dPLeE0Qm6lc2eFxxelP4LIzyX1QX/ioQRmwV6DZMe0BlcJPFjcX1yID6zRRYcxhH
ahpo6vGdleLZF44pFxRoFbctZY1YBCJW+gik4T9JxxrwGv0R+Cm0pKs2rYdfcjO4
mabNnfjGdygmOCi0YLAZ20rRgHgUr25cP5U+CCyvFwjdnsIbzgO22ebkI7bWk50V
vQHFCF04VIvRTrJ6KmFp+K8xW3Tsm9jCkEiJTrYQTtb7CpmeRlH23NUtfUW5qvN/
iQIDAQAB
-----END PUBLIC KEY-----";

        let err = pem_to_key(rsa_pem).unwrap_err();
        assert_eq!(err.code, InvalidPEM);
        assert_eq!(err.message, "Unexpected Crypto Algorithm");
    }
}
