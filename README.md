# Go As Your Backend

## What is this?

This is a starter template for using Go as a REST service. It uses [Fiber](https://docs.gofiber.io/) as the web framework and [GORM](https://gorm.io/) as the ORM.

## Why?

Because Go is awesome. And will serve more requests per dollar than most other languages.

## How to use?

After cloning this template repository & adding database credentials in a `.env` file in the root directory, the following Makefile commands are available:

- `make run-dev` - Runs the application in development mode
- `make build` - Builds the application
- `make run-prod` - Runs the application in production mode. This requires the application to be built first
- `make test` - Runs any Go test it finds
- `make clean` - Cleans the build files
- `make migrate` - Runs the database migrations

## Contributing

Feel free to open issues and pull requests. Any feedback is welcome!

## License

This project is licensed under the Unlicensed License - see the [UNLICENSE](UNLICENSE) file for details
