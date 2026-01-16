# Go As Your Backend

A starter template for building REST APIs in Go. Uses [Fiber](https://docs.gofiber.io/), [GORM](https://gorm.io/), and [Atlas](https://atlasgo.io/) for migrations, with a sensible project structure and built-in JWT authentication.

> **Note:** This template is extracted from a production backend, not designed as a comprehensive framework. The patterns here reflect real-world needs rather than theoretical completeness—which is exactly what makes it practical.

## Why?

Because Go is awesome, simple, reliable, and will serve more requests per dollar than most other languages.

## What does the template offer?

- JWT authentication with access/refresh token flow
- Database migrations with Atlas
- Generic repository pattern with dynamic filtering
- Role-based access control middleware
- Structured logging
- Project structure for scalable web applications

## Project Structure

```
app/
  logging/       # Logging utilities
  models/        # Database models and DTOs
  modules/       # Feature modules (users, posts, roles)
    users/       # Auth, handlers, service, repository
    posts/       # Example CRUD module
    roles/       # Role management
  spec/          # Generic repository interface and filters
  utils/         # Helper functions
```

## Dependencies

- [Fiber](https://gofiber.io/) - Fast HTTP framework built on Fasthttp
- [GORM](https://gorm.io/) - ORM for Go
- [Atlas](https://atlasgo.io/) - Database schema management and migrations
- [golang-jwt](https://github.com/golang-jwt/jwt) - JWT implementation

## Getting Started

1. Clone this repository
2. Copy `env-example` to `.env` and configure your database credentials
3. Generate the first migration:
```bash
make new-migration name=init
```
4. Run migrations:
```bash
make apply-migration
```
5. Start the server:
```bash
make build && ./bin/go-as-your-backend
```

## Environment Variables

See `env-example` for all required variables:

## API Endpoints

See `REQUESTS.md` for detailed endpoint documentation with curl examples.

### Auth

- `POST /api/v1/users/register` - Register new user
- `POST /api/v1/users/login` - Login
- `POST /api/v1/users/refresh` - Refresh access token
- `POST /api/v1/users/logout` - Logout
- `GET /api/v1/users/me` - Get current user

### Posts

- `GET /api/v1/posts` - List posts
- `GET /api/v1/posts/published` - List published posts
- `GET /api/v1/posts/count` - Get posts count
- `POST /api/v1/posts` - Create post (requires auth)

## Example Request

```bash
curl --get --data-urlencode 'filter={"where":{"and":[{"column":"published","operator":"=","value":true}]}}' http://localhost:3232/api/v1/posts
```

## Example Response

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

## Dynamic Query Filtering

The template includes a filter system for building dynamic queries from API parameters.

### Filter Syntax

```json
{
    "select": ["id", "title", "published"],
    "where": {
        "and": [{ "column": "published", "operator": "=", "value": true }],
        "or": [{ "column": "title", "operator": "LIKE", "value": "%hello%" }]
    },
    "order_by": [{ "column": "created_at", "direction": "DESC" }]
}
```

### Supported Operators

`=`, `>`, `<`, `>=`, `<=`, `LIKE`, `IN`, `IS NULL`

### Examples

Filter published posts:

```
/posts?filter={"where":{"and":[{"column":"published","operator":"=","value":true}]}}
```

With ordering:

```
/posts?filter={"where":{"and":[{"column":"published","operator":"=","value":true}]},"order_by":[{"column":"created_at","direction":"DESC"}]}
```

With pagination:

```
/posts?limit=10&offset=0
```

## Testing with Veriflow

This template includes a [Veriflow](https://github.com/okira-e/veriflow) configuration for API flow testing. Veriflow is a CLI tool for testing REST API flows with support for chained requests, assertions, and variable exports.

### Running Tests

```bash
veriflow run
```

### Test Configuration

See `veriflow.json` for the test flows:

- **auth** - Registration and login flow
- **posts** - CRUD operations for posts
- **filters** - Filter query parameter tests
- **validation** - Input validation tests

Example flow:

```json
{
    "name": "auth",
    "steps": [
        {
            "name": "register",
            "request": {
                "method": "POST",
                "path": "/users/register",
                "json": {
                    "email": "john.doe-{{RUN_ID}}@example.com",
                    "first_name": "John",
                    "last_name": "Doe",
                    "password": "$Password2025",
                    "phone": "+1234567890"
                }
            },
            "assert": {
                "status": 201
            }
        }
    ]
}
```

## Contributing

Feel free to suggest improvements and open a PR.

## License

This project is licensed under the Unlicense - see the [UNLICENSE](UNLICENSE) file for details.
