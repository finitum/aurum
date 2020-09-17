#[macro_use]
extern crate pest_derive;

use crate::parser::parse_line;
use aurum_rs::{Aurum, Role, User};
use itertools::Itertools;
use log::LevelFilter;
use std::env::args;
use std::io;
use std::io::{BufRead, Write};
use url::Url;

mod error;
mod parser;

pub enum Command {
    Connect(Url),
    Help,
    Login(String),
    Signup(String, String),
    Logout,
    User(Option<String>),
}

#[derive(Default)]
struct State {
    aurum: Option<Aurum>,
    user: Option<User>,
}

fn execute_line(line: String, state: &mut State) {
    if line.trim() == "" {
        return;
    }

    let command = match parse_line(line) {
        Ok(i) => i,
        Err(e) => {
            log::error!("Syntax error: {}", e);
            return;
        }
    };

    match command {
        Command::Connect(url) => match Aurum::new(url.to_string()) {
            Ok(i) => {
                log::info!("Succesfully connected to aurum");
                state.aurum = Some(i)
            }
            Err(e) => log::error!("Failed to connect: {:?}", e),
        },
        Command::Login(username) => {
            state.user = None;

            if let Some(aurum) = &mut state.aurum {
                let password = match rpassword::read_password_from_tty(Some("password>>")) {
                    Ok(i) => i,
                    Err(e) => {
                        log::error!("Failed to connect read password: {:?}", e);
                        return;
                    }
                };

                log::info!("logging in");

                match aurum.login(username.clone(), password) {
                    Ok(i) => {
                        log::info!("Succesfully logged in as {}", username);
                        state.user = Some(i)
                    }
                    Err(e) => log::error!("Failed to log in: {:?}", e),
                }
            } else {
                log::error!("Not connected to an aurum server. Try the `connect` command")
            }
        }
        Command::Signup(username, email) => {
            state.user = None;

            if let Some(aurum) = &mut state.aurum {
                let password = match rpassword::read_password_from_tty(Some("password>>")) {
                    Ok(i) => i,
                    Err(e) => {
                        log::error!("Failed to connect read password: {:?}", e);
                        return;
                    }
                };
                let password_repeat =
                    match rpassword::read_password_from_tty(Some("password repeat>>")) {
                        Ok(i) => i,
                        Err(e) => {
                            log::error!("Failed to connect read password: {:?}", e);
                            return;
                        }
                    };

                if password != password_repeat {
                    log::error!("Passwords don't match.");
                    return;
                }

                log::info!("Signing up");

                match aurum.signup(username.clone(), email, password) {
                    Ok(i) => {
                        log::info!("Succesfully signed up and logged in as {}", username);
                        state.user = Some(i)
                    }
                    Err(e) => log::error!("Failed to log in: {:?}", e),
                }
            } else {
                log::error!("Not connected to an aurum server. Try the `connect` command")
            }
        }
        Command::User(username) => {
            if let Some(aurum) = &mut state.aurum {
                let username = if let Some(username) = username {
                    username
                } else {
                    if let Some(user) = &state.user {
                        user.username().to_string()
                    } else {
                        log::error!("No username specified and not logged in");
                        return;
                    }
                };

                // aurum.user(username)
            }
        }
        Command::Logout => {
            state.user = None;
        }
        Command::Help => println!(
            "
Aurum command line interface

commands:
- help | ?                        Get this help page.
- connect   <url>                 Connect to an aurum server with this url.
- login     <username>            Login as this user on aurum. Prompts for a password.
- signup    <username> <email>    Register a new user. Prompts for a password.
- logout                          Logs you out. Alternatively you can log in as another user.
- user      [name]                Prints information about a user (without name uses yourself).

"
        ),
    }
}

fn main() {
    simple_logger::SimpleLogger::new()
        .with_level(LevelFilter::Info)
        .init()
        .expect("Failed to initialize logger");

    let mut state = State::default();

    let stdin = io::stdin();
    let mut stdout = io::stdout();

    let argstring = args().into_iter().skip(1).join(" ");
    log::debug!("arguments: {}", argstring);

    for i in argstring.split(",") {
        execute_line(i.to_string(), &mut state);
    }

    loop {
        if let Some(user) = &state.user {
            if user.role() == &Role::Admin {
                write!(stdout, "{} (admin) >>", user.username())
                    .expect("failed to write to stdout");
            } else {
                write!(stdout, "{} >>", user.username()).expect("failed to write to stdout");
            }
        } else {
            write!(stdout, ">>").expect("failed to write to stdout");
        }
        stdout.flush().expect("failed to flush stdout");

        if let Some(Ok(line)) = stdin.lock().lines().next() {
            execute_line(line, &mut state);
        } else {
            break;
        }
    }
}
