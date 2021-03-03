NAME ?= server

GO ?= go

.PHONY: all build clean test re

all: build

build:
	CGO_ENABLED=0 $(GO) build -o $(NAME) cmd/server/main.go

clean:
	rm -f $(NAME)

test: test-unit test-integration test-endtoend test-ui

test-unit:
	$(GO) test -v -p=1 -tags=unit ./...

test-integration:
	$(GO) test -v -p=1 -tags=integration ./...

test-endtoend:
	$(GO) test -v -p=1 -tags=endtoend ./...

test-ui:
	$(GO) test -v -p=1 -tags=ui ./...

re: clean all
