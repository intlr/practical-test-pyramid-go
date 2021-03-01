NAME ?= server

GO ?= go

.PHONY: all build clean test re

all: build

build:
	$(GO) build -o $(NAME) cmd/server/main.go

clean:
	rm -f $(NAME)

test:
	$(GO) test -v -p=1 ./...

re: clean all
