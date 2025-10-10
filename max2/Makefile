include hack/hack-cli.mk

.PHONY: all build run clean

all: build

build:
	@echo "Building..."
	go mod tidy
	go build -o bin/main main.go
	@echo "Build completed!"

run:
	@go run main.go

clean:
	@rm -rf bin/
	@echo "Clean completed!"

help:
	@echo "Available targets:"
	@echo "  build - Build the application"
	@echo "  run   - Run the application"
	@echo "  clean - Clean build artifacts"