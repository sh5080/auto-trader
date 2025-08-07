run:
	PORT=8087 go run cmd/trader/main.go

swagger:
	./scripts/generate_swagger.sh