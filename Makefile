dev:
	nodemon --watch './**/*.go' --signal SIGTERM --exec 'go' run ./main.go

build-exe:
	go build -o ./build/server ./main.go

prod:
	# Make sure to run `make build-exe` before running this command
	./build/server

test:
	go test ./src/... -v

migrate:
	go run ./main.go migrate
