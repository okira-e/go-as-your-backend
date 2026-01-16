# Load environment variables
include .env

# === Commands ===
APP := go-as-your-backend

build:
	go build -o bin/$(APP) .

release:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/$(APP)
	
new-migration:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make new-migration name=<migration_name>"; \
		exit 1; \
	fi
	atlas migrate diff $(name) \
		--dir file://migrations \
		--to file://schema.hcl \
		--dev-url "docker://postgres/15/dev?search_path=public"

apply-migration:
	atlas migrate apply \
		--dir file://migrations \
		--url "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable&search_path=public"
