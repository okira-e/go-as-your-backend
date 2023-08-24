# Go As Your Backend

## What is this?

This is a starter template for using Go as a REST service. It uses [Fiber](https://docs.gofiber.io/) as the web framework and [GORM](https://gorm.io/) as the ORM.

## Why?

Because Go is awesome. And will serve more requests per dollar than most other languages.

## What's does the template offer?

- A Makefile with commands for running the application in development and production mode, building the application, running tests and running database migrations
- Swagger documentation for the API
- Project structure for a web application
- Web framework and ORM setup

## Usage

After cloning this template repository & adding database credentials in a `.env` file in the root directory, the following Makefile commands are available:

- `make dev` - Runs the application in development mode
- `make build-exe` - Builds the application
- `make prod` - Runs the application in production mode. This requires the application to be built first
- `make test` - Runs any Go test it finds
- `make migrate-db` - Runs the database migrations
- `make openapi` - Generates/Updates the OpenAPI documentation

## Contributing

Feel free to open issues and pull requests. Any feedback is welcome!

## License

This project is licensed under the Unlicensed License - see the [UNLICENSE](UNLICENSE) file for details
