---
permalink: /
---
Personal notes on [the article Ham Vocke wrote](https://martinfowler.com/articles/practical-test-pyramid.html) on the subject and interpretation to a Golang application

## Unit tests

- Foundation of the test suite
- Make sure a _System Under Test (SUT)_ works as intended
- Unit tests are fast and so more common that other types of tests
- Solitary unit tests are tests doubling all collaborators
- Sociable unit tests are tests allowing communications with real
  collaborators
- _Test-Driven Development (TDD)_ lets unit tests guide development
- One test class per production class rule of thumb
- Unit test at least public interfaces
- Includes happy and edge cases, without being too tied to implementation
- Arrange, Act, Assert

```go
package main

import (
	"testing"

	"github.com/alr-lab/test-double-go/service"
)

const email = "fake"

type StubStore struct{}

func (s StubStore) GetCustomerEmail(id int) string {
	return email
}

func TestService_Get(t *testing.T) {
	serv := service.New(StubStore{})

	got := serv.Get()
	if got != email {
		t.Fatalf("got %q, want %q", got, email)
	}
}
```

## Integration tests

## Contract tests

## UI tests

## End-to-end tests

## Acceptance tests

## Exploratory testing

## Terminology

## Deployment pipeline

## Test duplication

## Clean test code
