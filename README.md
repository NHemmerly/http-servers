# http-servers
Small implementation of a RESTful-like API for a fictional web app called "Chirpy". Chirpy implements users and chirps that belong to those users. 

# API Documentation

## User Resource
```json
{
    "id": "<user uuid>",
    "created_at": "<creation timestamp>",
    "updated_at": "<timestamp of last update>",
    "email": "<user@email.com>",
    "is_chirpy_red": false
}
```

### POST /api/users
This endpoint can be used to create a new user. It accepts a body:
```json
{
    "email":"<newuser@email.com>",
    "password":"<new password>"
}
```

Response:
```json
{
    "id":"<uuid>",
    "created_at": "<creation timestamp>",
    "updated_at": "<timestamp of last update>",
    "email":"<newuser@email.com>",
    "is_chirpy_red": false
}
```

### POST /api/login
Logs in a user and provides them with a new access token and refresh token.

Request:
```json
{
    "email":"<newuser@email.com>",
    "password":"<new password>"
}
```

Response:
```json
{
    "id":"<uuid>",
    "created_at": "<creation timestamp>",
    "updated_at": "<timestamp of last update>",
    "email":"<newuser@email.com>",
    "token":"<access_token_string>",
    "refresh_token":"<refresh token string>",
    "is_chirpy_red": false
}
```

### PUT /api/users
Updates a user's login with information provided in the request.

Request:
```json
{
    "token":"<access token string>",
    "password":"<new password>"
}
```

Response:
```json
{
    "id":"<uuid>",
    "created_at": "<creation timestamp>",
    "updated_at": "<timestamp of last update>",
    "email":"<user@email.com>",
    "is_chirpy_red": false
}
```

## Chirp Resource
```json
{
    "id":"<uuid>",
    "created_at": "<creation timestamp>",
    "updated_at": "<timestamp of last update>",
    "body":"<chirp content>",
    "user_id":"<uuid of chirp author>"
}
```

### POST /api/chirps
Posts a chirp as the currently authenticated user. 

Request:
```json
Header:
{
    "Authorization": "Bearer <token>"
}
Body:
{
    "id":"<uuid>",
    "created_at": "<creation timestamp>",
    "updated_at": "<timestamp of last update>",
    "body":"<chirp content>",
    "user_id":"<uuid of chirp author>"
}
```

