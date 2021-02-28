NAME = a.out

GO ?= go

.PHONY: all build clean test re

all: build

build:
	CGO_ENABLED=0 GOOS=linux $(GO) build -o $(NAME) cmd/server/main.go

clean:
	rm -f $(NAME)

test:
	$(GO) test -v ./...

re: clean all
