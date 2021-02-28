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
// Package service provides a simple service on which we can experiment
// tests
//
// Part of the Test Double repository I also published
// https://github.com/alr-lab/test-double-go
package service

type (
	// Service describes a service
	Service struct {
		store Store
	}

	// Store defines a contract for a datastore
	Store interface {
		GetCustomerEmail(id int) string
	}
)

// New returns a new service
func New(store Store) Service {
	return Service{store: store}
}

// Get returns a specific customer email
func (s Service) Get() string {
	return s.store.GetCustomerEmail(42)
}
```

```go
// Simple unit test
//
// Part of the Test Double repository I also published
// https://github.com/alr-lab/test-double-go
//
// The System Under Test is the ``Service'' object, and both the public
// interfaces ``New'' and ``Get'' are tested. It is considered a solitary
// unit test as we are doubling the datastore.
package service_test

import (
	"testing"

	"github.com/alr-lab/test-double-go/service"
)

const email = "fake"

type StubStore struct{}

func (s StubStore) GetCustomerEmail(id int) string {
	return email
}

func TestService(t *testing.T) {
	// Arrange
	serv := service.New(StubStore{})

	// Act
	got := serv.Get()

	// Assert
	if got != email {
		t.Fatalf("got %q, want %q", got, email)
	}
}
```

## Integration tests

- Test integrations with other parts such as database, filesystems, and
  network
- Also run the components we are integrating
- Different kind of integration tests
  - Test through the entire stack
  - Test integrations one by one, doubling the others if needed
- Database integration test
  1. Start database
  2. Connect application to database
  3. Interact with database
  4. Validate expectations
- API integration test
  1. Start application
  2. Start instance of the API
  3. Interact with the API
  4. Validate expectations
- Integration tests are slower than unit-tests doubling integrations

```go
package dbtesting

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func DatabaseHelper(t *testing.T) *sql.DB {
	t.Helper()

	// Open MySQL connection and fill database with testing data...

	return conn
}
```

```go
package store

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct{
	conn *sql.DB
}

func (st *Store) SetConn(conn *sql.DB) {
	st.conn = conn
}

func (s Store) GetCustomerEmail(ctx context.Context, id int) (string, error) {
	query := `SELECT email FROM Customers WHERE id = ?`

	args := []interface{}{id}
	res, err := s.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return "", fmt.Errorf("unable to query database, err = %s", err)
	}
	defer func() { _ = res.Close() }()

	email := ""
	if res.Next() {
		if err := res.Scan(&email); err != nil {
			return email, fmt.Errorf("unable to parse customer email, err = %s", err)
		}
	}

	return email, nil
}
```

```go
package store_test

import (
	"context"
	"testing"

	"github.com/alr-lab/ptp/internal/dbtesting"
	"github.com/alr-lab/ptp/store"
)

func TestStore(t *testing.T) {
	want := "fake"

	conn := dbtesting.DatabaseHelper(t)
	defer func(){ _ = conn.Close() }()

	st := &store.Store{}
	st.SetConn(conn)

	got, err := st.GetCustomerEmail(context.Background(), 42)
	if err != nil {
		t.Fatalf("Unable to get customer email, err = %s", err)
	}

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
```

## Contract tests

## UI tests

## End-to-end tests

## Acceptance tests

## Exploratory testing

## Terminology

## Deployment pipeline

## Test duplication

## Clean test code
