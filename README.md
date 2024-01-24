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

## Dependencies

- [swag](https://github.com/swaggo/swag): To generate OpenAPI documentation using `make docs` you need to install [swag](https://github.com/swaggo/swag) system-wide. Installation is detailed in their Github.
- [Go](https://go.dev/): duh

## Usage

After cloning this template repository & adding database credentials in a `.env` file in the root directory, the following Makefile commands are available:

- `make dev` - Runs the application in development mode
- `make build` - Builds the application
- `make prod` - Runs the application in production mode. This requires the application to be built first
- `make test` - Runs any Go test it finds
- `make migrate-db` - Runs the database migrations
- `make docs` - Generates/Updates the OpenAPI documentation. Check the Note section for more details.

## Note

Generation OpenAPI documentation is seamless. However, for some reason with the technologies that I am using to generate and view the /swagger/index.html/ page there seems to be caching issues where Newly generated documents are not shown in the browser. If you ever encounter an issue with the content not being updated after running `make docs` just clear the cache or view the page from another browser window (anonymous window is optimal.)

## Contributing

Feel free to open issues and pull requests. Any feedback is welcome!

## License

This project is licensed under the Unlicensed License - see the [UNLICENSE](UNLICENSE) file for details
