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

Response 201 Created:
```json
{
    "id":"<chirp id>",
    "created_at": "<creation timestamp>",
    "updated_at": "<timestamp of last update>",
    "body":"<chirp content>",
    "user_id":"<uuid of chirp author>"
}
```

### GET /api/chirps?{author_id=uuid&sort=asc|desc}
Returns a set of chirps depending on whether a user id was provided as a query parameter. Will also optionally sort the chirps based on the "sort" query parameter. If no user id is provided, the request will return all chirps in ascending order. 

Response 200 OK:
```json
{
    [
        {
            "id":"<chirp id>",
            "created_at": "<creation timestamp>",
            "updated_at": "<timestamp of last update>",
            "body":"<chirp content>",
            "user_id":"<uuid of chirp author>"
        },
        {
            "id":"<chirp id>",
            "created_at": "<creation timestamp>",
            "updated_at": "<timestamp of last update>",
            "body":"<chirp content>",
            "user_id":"<uuid of chirp author>"
        },
    ]
}
```

### GET /api/chirps/{chirp_id}
Returns a single chirp based on a unique chirp ID.

Response 200 OK:
```json
{
    "id":"<chirp id>",
    "created_at": "<creation timestamp>",
    "updated_at": "<timestamp of last update>",
    "body":"<chirp content>",
    "user_id":"<uuid of chirp author>"
}
```

### DELETE /api/chirps/{chirp_id}
Deletes the authorized user's chirp after validating that it belongs to them. Request must include an access token in the header and a chirp ID in the request path.

Request:
```json
Header:
{
    "Authorization": "Bearer <token>",
}
```

Response 204 No Content:
>"Chirp deleted"

## Auth Endpoints
### POST /api/refresh
Requires a refresh token in the header. Replies with a new access token.

Response:
```json
{
    "token":"<new access token>"
}
```

### POST /api/revoke
Requires a refresh token in the header. Revokes and makes the request's refresh token invalid.

Response 204 No Content

## Admin Endpoints
### POST /admin/reset
If the requesting client has all of the necessary environment variables, the backend database will be fully cleared of users and chirps.

### GET /admin/metrics
Returns the total number of "hits" on the application's user-facing endpoints. 

