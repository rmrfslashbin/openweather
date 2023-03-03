.DEFAULT_GOAL := default

build:
	@echo "Building..."
	@if [ ! -d "./bin" ]; then mkdir bin; fi
	@go build -o bin/openweather cmd/main.go

install:
	@go install

tidy:
	@echo "Making mod tidy"
	@go mod tidy

update:
	@echo "Updating..."
	@go get -u ./...
	@go mod tidy

default: tidy build
