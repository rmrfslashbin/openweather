.DEFAULT_GOAL := build
.PHONY: build

GIT_HASH := $(shell git rev-parse HEAD | cut -c 1-8)

build:
	@echo build hash is $(GIT_HASH)
	@printf "  building openweather:\n"
	@printf "    linux  :: arm64"
	@GOOS=linux GOARCH=arm64 go build -ldflags "-X github.com/rmrfslashbin/openweather/cmd.Version=$(BUILD_DATE)" -o bin/openweather-linux-arm64 .
	@printf " done.\n"
	@printf "    linux  :: amd64"
	@GOOS=linux GOARCH=amd64 go build -ldflags "-X github.com/rmrfslashbin/openweather/cmd.Version=$(BUILD_DATE)" -o bin/openweather-linux-amd64 .
	@printf " done.\n"
	@printf "    darwin :: amd64"
	@GOOS=darwin GOARCH=amd64 go build -ldflags "-X github.com/rmrfslashbin/openweather/cmd.Version=$(BUILD_DATE)" -o bin/openweather-darwin-amd64 .
	@printf " done.\n"
	@printf "    darwin :: arm64"
	@GOOS=darwin GOARCH=arm64 go build -ldflags "-X github.com/rmrfslashbin/openweather/cmd.Version=$(BUILD_DATE)" -o bin/openweather-darwin-arm64 .
	@printf " done.\n"