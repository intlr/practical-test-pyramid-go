# Practical Test Pyramid with Go

[![GoDoc](https://godoc.org/github.com/alr-lab/practical-test-pyramid-go?status.svg)](https://godoc.org/github.com/alr-lab/practical-test-pyramid-go)
[![Automated tests](https://github.com/alr-lab/practical-test-pyramid-go/actions/workflows/push.yml/badge.svg)](https://github.com/alr-lab/practical-test-pyramid-go/actions/workflows/push.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/alr-lab/practical-test-pyramid-go)](https://goreportcard.com/report/github.com/alr-lab/practical-test-pyramid-go)

Practical Test Pyramid with Go is an experiment to work with
all layers of the [_Practical Test Pyramid_][ptp-ham] introduced by Ham
Vocke in 2018.

## Getting started

```
$ docker-compose up -d --build
```

## Run tests

```
$ docker exec -it app make test
```

[ptp-ham]: https://martinfowler.com/articles/practical-test-pyramid.html
