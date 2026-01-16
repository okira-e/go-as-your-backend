# Endpoints

## Auth Flow

### Registration

- Request

```sh
curl -X POST http://localhost:3232/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "password": "$Password2025",
    "phone": "+1234567890"
  }'
```

- Response:

```json
{
    "success": true,
    "status": 201,
    "message": "User created successfully",
    "data": {
        "id": "019a8880-d7f8-7e72-81cb-75ae25e651e9",
        "role_id": null,
        "first_name": "John",
        "last_name": "Doe",
        "email": "john.doe@example.com",
        "created_at": "2025-11-15T17:12:16.633114Z",
        "updated_at": "2025-11-15T20:12:16.634604+03:00"
    }
}
```

### Login

- Request

```sh
curl -X POST http://localhost:3232/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "$Password2025"
  }' \
  -c cookies.txt
```

- Response:

```json
{
    "success": true,
    "status": 200,
    "message": "Login successful",
    "data": {
        "id": "019a529c-c734-7796-ba61-81fe04e75647",
        "role_id": null,
        "first_name": "John",
        "last_name": "Doe",
        "email": "john.doe@example.com",
        "created_at": "2025-11-05T06:03:17.684661Z",
        "updated_at": "2025-11-05T09:03:17.684893Z"
    }
}
```

### Me

- Request

```sh
curl -X GET http://localhost:3232/api/v1/users/me \
  -H "Content-Type: application/json" \
  -b cookies.txt
```

- Response:

```json
{
    "success": true,
    "status": 200,
    "message": "",
    "data": {
        "id": "019a529c-c734-7796-ba61-81fe04e75647",
        "role_id": null,
        "first_name": "John",
        "last_name": "Doe",
        "email": "john.doe@example.com",
        "created_at": "2025-11-05T06:03:17.684661Z",
        "updated_at": "2025-11-05T09:03:17.684893Z"
    }
}
```

### Refresh Token

- Request

```sh
curl -X POST http://localhost:3232/api/v1/users/refresh -b cookies.txt
```

- Response:

```json
{
    "success": true,
    "status": 200,
    "message": "Token refreshed successfully",
    "data": {}
}
```

### Logout

- Request

```sh
curl -X POST http://localhost:3232/api/v1/users/logout -b cookies.txt
```

- Response:

```json
{ "success": true, "status": 200, "message": "Logout successful", "data": null }
```

### User Contact Info

- Request

```sh
curl http://localhost:3232/api/v1/users/contact-info/${user_id}
```

- Response:

```json
{
    "success": true,
    "status": 200,
    "message": "",
    "data": { "phone": "+1234567890" }
}
```

## Posts Flow

### Create Post

- Request

```sh
curl -X POST http://localhost:3232/api/v1/posts/ \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "title": "My First Post",
    "content": "This is the content of my first post.",
    "published": true
  }'
```

- Response

```json
{
    "success": true,
    "status": 201,
    "message": "",
    "data": {
        "id": "019a86ba-fc66-7150-8e23-307b5db2c5e9",
        "title": "My First Post",
        "content": "This is the content of my first post.",
        "published": true,
        "user_id": "019a86ad-0e55-79a6-b314-74e5d8a06848",
        "created_at": "2025-11-15T17:12:16.633114Z",
        "updated_at": null
    }
}
```

### Get Posts Count

- Request

```sh
curl http://localhost:3232/api/v1/posts/count/
```

- Response

```json
{ "success": true, "status": 200, "message": "", "data": 8 }
```

### List Posts

- Request

```sh
curl http://localhost:3232/api/v1/posts/
```

- Response

```json
{
    "success": true,
    "status": 200,
    "message": "",
    "data": [
        {
            "id": "019a8318-66eb-7824-89c6-c9bde9ea9cbe",
            "title": "My First Post",
            "content": "This is the content of my first post.",
            "published": true,
            "user_id": "019a529c-c734-7796-ba61-81fe04e75647",
            "created_at": "2025-11-15T17:12:16.633114Z",
            "updated_at": null
        },
        {
            "id": "019a8318-822b-7ecf-865a-3ec8dfe60aea",
            "title": "Draft Post",
            "content": "This is a draft.",
            "published": false,
            "user_id": "019a529c-c734-7796-ba61-81fe04e75647",
            "created_at": "2025-11-15T18:00:00.000000Z",
            "updated_at": null
        }
    ]
}
```

### Published Posts

- Request

```sh
curl http://localhost:3232/api/v1/posts/published \
  -H "Content-Type: application/json"
```

- Response

```json
{
    "success": true,
    "status": 200,
    "message": "",
    "data": [
        {
            "id": "019a8318-66eb-7824-89c6-c9bde9ea9cbe",
            "title": "My First Post",
            "content": "This is the content of my first post.",
            "published": true,
            "user_id": "019a529c-c734-7796-ba61-81fe04e75647",
            "created_at": "2025-11-15T17:12:16.633114Z",
            "updated_at": null
        }
    ]
}
```

## Filtering

The API supports a flexible filter query parameter for filtering, sorting, and selecting fields.

### Filter by Title

```sh
curl "http://localhost:3232/api/v1/posts?filter={\"where\":{\"and\":[{\"column\":\"title\",\"operator\":\"=\",\"value\":\"My First Post\"}]}}"
```

### Filter Published Posts

```sh
curl "http://localhost:3232/api/v1/posts?filter={\"where\":{\"and\":[{\"column\":\"published\",\"operator\":\"=\",\"value\":true}]}}"
```

### Order By with Filter

```sh
curl "http://localhost:3232/api/v1/posts?filter={\"where\":{\"and\":[{\"column\":\"published\",\"operator\":\"=\",\"value\":true}]},\"order_by\":[{\"column\":\"created_at\",\"direction\":\"DESC\"}]}"
```

### Select Specific Fields

```sh
curl "http://localhost:3232/api/v1/posts?filter={\"select\":[\"id\",\"title\",\"published\"]}"
```

### Pagination

```sh
curl "http://localhost:3232/api/v1/posts?limit=10&offset=0"
```
