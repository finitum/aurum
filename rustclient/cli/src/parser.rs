use pest::iterators::Pair;
use crate::error::CliError;
use crate::Command;
use url::Url;
use pest::Parser;

#[derive(Parser)]
#[grammar = "language.pest"]
struct CommandParser;


fn parse_connect<'a>(mut pairs: impl Iterator<Item = Pair<'a, Rule>>) -> Result<Command, CliError>  {
    let string = pairs.next().ok_or(CliError::Custom("failed to parse command".to_string()))?.as_str();

    let string = if !string.starts_with("http://") & !string.starts_with("https://")  {
        format!("http://{}", string)
    } else {
        string.to_string()
    };

    let url = Url::parse(&string)
        .map_err(|e| CliError::Custom(format!("failed to parse url: {}", e)))?;

    Ok(Command::Connect(url))
}

fn parse_login<'a>(mut pairs: impl Iterator<Item = Pair<'a, Rule>>) -> Result<Command, CliError>  {
    let username = pairs.next().ok_or(CliError::Custom("failed to parse command".to_string()))?.as_str();
    Ok(Command::Login(username.to_string()))
}

fn parse_signup<'a>(mut pairs: impl Iterator<Item = Pair<'a, Rule>>) -> Result<Command, CliError>  {
    let username = pairs.next().ok_or(CliError::Custom("failed to parse command".to_string()))?.as_str();
    let email = pairs.next().ok_or(CliError::Custom("failed to parse command".to_string()))?.as_str();
    Ok(Command::Signup(username.to_string(), email.to_string()))
}

fn parse_user<'a>(mut pairs: impl Iterator<Item = Pair<'a, Rule>>) -> Result<Command, CliError>  {
    let username = if let Some(nxt) = pairs.next() {
        Some(nxt.as_str().to_string())
    } else {
        None
    };

    Ok(Command::User(username))
}

pub fn parse_line(line: String) -> Result<Command, CliError> {
    let mut pairs = CommandParser::parse(Rule::line, &line)?.into_iter();

    let first = pairs.next().ok_or(CliError::Custom("failed to parse command".to_string()))?;

    let command = match first.as_rule() {
        Rule::connect => parse_connect(first.into_inner())?,
        Rule::help => Command::Help,
        Rule::login => parse_login(first.into_inner())?,
        Rule::signup => parse_signup(first.into_inner())?,
        Rule::logout => Command::Logout,
        Rule::user => parse_user(first.into_inner())?,
        _ => todo!(),
    };

    Ok(command)
}
