run:
	PORT=8087 go run cmd/trader/main.go

build:
	make lint
	go build ./...

swagger:
	./scripts/generate_swagger.sh

migrate:
	./scripts/migrate.sh

lint:
	@golangci-lint run ./...