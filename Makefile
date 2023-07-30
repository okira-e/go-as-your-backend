dev:
	nodemon --watch './**/*.go' --signal SIGTERM --exec 'go' run ./src/main.go

build:
	mkdir -p ./build
	go build -o ./build/main ./src/main.go

prod:
	# Make sure to run `make build` before running this command
	./build/main

test:
	go test ./src/... -v

clean:
	rm -rf ./build

migrate:
	go run ./src/main.go migrate
