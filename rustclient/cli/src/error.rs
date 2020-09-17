use std::error::Error;
use std::fmt;

#[derive(Debug)]
pub enum CliError {
    Pest(pest::error::Error<crate::parser::Rule>),
    Custom(String),
}

impl Error for CliError {}
impl fmt::Display for CliError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            CliError::Pest(e) => write!(f, "{}", e),
            CliError::Custom(e) => write!(f, "{}", e),
        }
    }
}
impl From<pest::error::Error<crate::parser::Rule>> for CliError {
    fn from(e: pest::error::Error<crate::parser::Rule>) -> Self {
        CliError::Pest(e)
    }
}
