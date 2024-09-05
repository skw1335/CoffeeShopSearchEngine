# Makefile for Go project

GO_FILES := $(wildcard *.go)
PKG_NAME = Go_Backend 
BIN_NAME = your_binary_name

run:
	go run main.go

build:
	@echo "Building..."
	go build -o ${BIN_NAME}

test:
	@echo "Running tests..."
	go test -v ./...

migration:
	@migrate create -ext sql -dir migrate/migrations

migrate-up:
	go run migrate/main.go up 

migrate-down:
	go run migrate/main.go down


