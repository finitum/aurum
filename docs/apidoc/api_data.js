define({ "api": [
  {
    "type": "post",
    "url": "/login",
    "title": "Login",
    "description": "<p>Logs a user in returning a tokenpair</p>",
    "name": "Login",
    "group": "Authentication",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "username",
            "description": "<p>The user's username</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "password",
            "description": "<p>The user's password</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "Request Example:",
          "content": "{\n\t\"username\": \"victor\",\n\t\"password\": \"hunter2\"\n}",
          "type": "json"
        }
      ]
    },
    "success": {
      "fields": {
        "Success 200": [
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "login_token",
            "description": "<p>The user's login token</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "refresh_token",
            "description": "<p>The user's refresh token</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "Success Response:",
          "content": "{\n\t\"login_token\": \"<JWT Token here>\"\n\t\"refresh_token\": \"<JWT Token here>\"\n}",
          "type": "json"
        }
      ]
    },
    "error": {
      "fields": {
        "Error 4xx": [
          {
            "group": "Error 4xx",
            "optional": false,
            "field": "400",
            "description": "<p>If an invalid body is provided.</p>"
          },
          {
            "group": "Error 4xx",
            "optional": false,
            "field": "401",
            "description": "<p>If the user does not exist or the password is wrong</p>"
          }
        ]
      }
    },
    "version": "0.0.0",
    "filename": "../core/web/AuthEndpoints.go",
    "groupTitle": "Authentication"
  },
  {
    "type": "post",
    "url": "/refresh",
    "title": "Refresh Token",
    "description": "<p>Refreshes your login token by using your refresh token</p>",
    "name": "Refresh",
    "group": "Authentication",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "refresh_token",
            "description": "<p>The refresh token to use.</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "Request Example:",
          "content": "{\n\t\"refresh_token\": \"<JWT Token here>\"\n}",
          "type": "json"
        }
      ]
    },
    "success": {
      "fields": {
        "Success 200": [
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "login_token",
            "description": "<p>A renewed login token</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "Success Response:",
          "content": "{\n\t\"login_token\": \"<JWT Token here>\"\n}",
          "type": "json"
        }
      ]
    },
    "error": {
      "fields": {
        "Error 4xx": [
          {
            "group": "Error 4xx",
            "optional": false,
            "field": "400",
            "description": "<p>If an invalid body or token is provided</p>"
          },
          {
            "group": "Error 4xx",
            "optional": false,
            "field": "404",
            "description": "<p>If the user does not exist (anymore)</p>"
          }
        ]
      }
    },
    "version": "0.0.0",
    "filename": "../core/web/AuthEndpoints.go",
    "groupTitle": "Authentication"
  },
  {
    "type": "post",
    "url": "/signup",
    "title": "Register",
    "description": "<p>Creates a new account</p>",
    "name": "Signup",
    "group": "Authentication",
    "parameter": {
      "fields": {
        "Parameter": [
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "username",
            "description": "<p>The username of the user</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "password",
            "description": "<p>The password of the user</p>"
          },
          {
            "group": "Parameter",
            "type": "String",
            "optional": false,
            "field": "email",
            "description": "<p>The E-Mail of the user</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "Request Example:",
          "content": "{\n\t\"username\": \"victor\",\n\t\"password\": \"hunter2\",\n\t\"email\": \"victor@example.com\"\n}",
          "type": "json"
        }
      ]
    },
    "success": {
      "examples": [
        {
          "title": "Success Response:",
          "content": "HTTP/1.1 201 Created",
          "type": "String"
        }
      ]
    },
    "error": {
      "fields": {
        "Error 4xx": [
          {
            "group": "Error 4xx",
            "optional": false,
            "field": "400",
            "description": "<p>If an invalid body is provided</p>"
          }
        ]
      }
    },
    "version": "0.0.0",
    "filename": "../core/web/AuthEndpoints.go",
    "groupTitle": "Authentication"
  },
  {
    "type": "get",
    "url": "/me",
    "title": "Request user info",
    "name": "GetUser",
    "group": "User",
    "header": {
      "fields": {
        "Authorization": [
          {
            "group": "Authorization",
            "type": "String",
            "optional": false,
            "field": "Authorization",
            "description": "<p>Users JWT Token</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "Authorization Example:",
          "content": "Authorization: \"Bearer <token>\"",
          "type": "String"
        }
      ]
    },
    "success": {
      "fields": {
        "Success 200": [
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "username",
            "description": "<p>The username of the user</p>"
          },
          {
            "group": "Success 200",
            "type": "String",
            "optional": false,
            "field": "email",
            "description": "<p>The E-Mail of the user</p>"
          },
          {
            "group": "Success 200",
            "type": "Number",
            "optional": false,
            "field": "role",
            "description": "<p>The role of the user (0 = User, 1 = Admin)</p>"
          },
          {
            "group": "Success 200",
            "type": "Boolean",
            "optional": false,
            "field": "blocked",
            "description": "<p>If the user is blocked</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "Success Response:",
          "content": "{\n\t\"username\":\"victor\",\n\t\"email\":\"victor@example.com\",\n\t\"role\":0,\n\t\"blocked\": false\n}",
          "type": "json"
        }
      ]
    },
    "error": {
      "fields": {
        "Error 4xx": [
          {
            "group": "Error 4xx",
            "optional": false,
            "field": "404",
            "description": "<p>If the user does not exist (anymore).</p>"
          }
        ]
      }
    },
    "version": "0.0.0",
    "filename": "../core/web/UserEndpoints.go",
    "groupTitle": "User"
  },
  {
    "type": "get",
    "url": "/users",
    "title": "Get all users",
    "description": "<p>Gets all users currently registered</p>",
    "name": "GetUsers",
    "group": "User",
    "permission": [
      {
        "name": "admin",
        "title": "Admin user",
        "description": "<p>Only available to admins, the first user of the server is by default admin.</p>"
      }
    ],
    "header": {
      "fields": {
        "Authorization": [
          {
            "group": "Authorization",
            "type": "String",
            "optional": false,
            "field": "Authorization",
            "description": "<p>Users JWT Token</p>"
          }
        ]
      },
      "examples": [
        {
          "title": "Authorization Example:",
          "content": "Authorization: \"Bearer <token>\"",
          "type": "String"
        }
      ]
    },
    "version": "0.0.0",
    "filename": "../core/web/UserEndpoints.go",
    "groupTitle": "User"
  }
] });
