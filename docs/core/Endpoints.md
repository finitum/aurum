# Endpoints

## Authentication
An authenticated request can be made by setting the `Authorization` HTTP Header to: `Bearer <token>`. Where `<token>` is
your jwt login token.

This is needed for all requests marked `Authenticated`.

### /signup - POST
Use the signup endpoint to create a new account.
The request should contain a JSON body with account information.

```json5
{
    "username": "xxx",      // Required
    "password": "yyy",      // Required
    "email": "zzz@xxx.com", // Required
}
```


Note: The first user to sign up will automatically become an administrator. 
Admininstrator users can make other users administrator.

### /login - POST

Use the login endpoint to log into your account and receive a refresh and login token. 
The request should contain a JSON body with login information.

An example request body would be:
```json5
{
    "username": "xxx",  // Required
    "password": "yyy",  // Required
}
```

which would return a response in the form of:
```json5
{
     "login_token": "JWT Login Token",
     "refresh_token": "JWT Refresh Token" 
}
```

**Note:** Be careful with the storage and handling of refresh tokens as the validity on them is high.

### /refresh - POST
Request a new login token by providing a refresh token.
```json5
{
    "refresh_token": "JWT Refresh Token", // Required
}
```

which would return a response in the form of:
```json5
{
     "login_token": "JWT Login Token"
}
```

### /revoke/ - DELETE  (TODO) 

Revoke *all* your tokens


## User

### /me - GET - Authenticated

Gets all your own user information.

this will return a response containing:
```json5
{
  "username": "my name",
  "email" : "name@mail.me",
  "role":  0,
  "blocked": false
}
```


### /me - PUT - Authenticated (TODO)
Update your user info

```json5
{
    "password": "yyy",      // Optional
    "email": "zzz@xxx.com", // Optional
}
```

Note: You can't change your username.

## Admin

