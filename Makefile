.PHONY: run
run:
	go run cmd/auth

.PHONY: build
build:
	go build -o bot cmd/auth