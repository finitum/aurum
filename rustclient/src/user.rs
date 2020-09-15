use serde::{Serialize, Deserialize};

#[repr(u8)]
#[derive(Serialize, Deserialize)]
pub(crate) enum Role {
    User = 0,
    Admin = 1,
}

impl Default for Role {
    fn default() -> Self {
        Role::User
    }
}

#[derive(Serialize, Deserialize, Default)]
pub(crate) struct User {
    pub(crate) username: String,
    pub(crate) password: String,
    pub(crate) email: String,
    pub(crate) role: Role,
    pub(crate) blocked: bool,
}
