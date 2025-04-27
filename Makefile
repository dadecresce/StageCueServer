APP_NAME := StageCueServer
VERSION  ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS  := -X 'main.version=$(VERSION)' -s -w

.PHONY: build build-win build-all test run clean

build:
	GOOS=$(shell go env GOOS) GOARCH=$(shell go env GOARCH) \
	go build -ldflags "$(LDFLAGS)" -o bin/$(APP_NAME) ./cmd/server

build-win:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 \
	go build -ldflags "$(LDFLAGS)" -o bin/$(APP_NAME).exe ./cmd/server

build-all: build build-win

test:
	go test ./... -cover

run: build
	./bin/$(APP_NAME) -config config.sample.toml

clean:
	rm -rf bin
