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

import "context"

type (
	// Service describes a service
	Service struct {
		store Store
	}

	// Store defines a contract for a datastore
	Store interface {
		GetCustomerEmail(ctx context.Context, id int) string
	}
)

// New returns a new service
func New(store Store) Service {
	return Service{store: store}
}

// Get returns a specific customer email
func (s Service) Get(ctx context.Context) string {
	return s.store.GetCustomerEmail(ctx, 42)
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
	"context"
	"testing"

	"github.com/alr-lab/test-double-go/service"
)

const email = "fake"

type StubStore struct{}

func (s StubStore) GetCustomerEmail(_ context.Context, _ int) string {
	return email
}

func TestService(t *testing.T) {
	// Arrange
	serv := service.New(StubStore{})
	ctx := context.Background()

	// Act
	got := serv.Get(ctx)

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

### Database integration

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
	// Arrange
	conn := dbtesting.DatabaseHelper(t)
	defer func() { _ = conn.Close() }()
	st := &store.Store{}
	st.SetConn(conn)
	tt := map[string]struct {
		id   int
		want string
	}{
		"Valid customer identifier will return valid email": {
			id:   42,
			want: "fake",
		},
		"Invalid customer identifier will return empty string": {
			id:   1337,
			want: "",
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			// Act
			got, err := st.GetCustomerEmail(context.Background(), tc.id)

			// Assert
			if err != nil {
				t.Fatalf("Unable to get customer email, err = %s", err)
			}
			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}
```

### External API integration

```go
// The ``mockapi'' application is a simple server acting as a replacement
// for an external API.
//
// With a ``docker-compose'' file we can start the server and allow our
// application to test its integration with a production-ready instance.
package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, _ *http.Request) {
		log.Print("Serving request")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		io.WriteString(w, `{"message":"Hello, world!"}`)
	})

	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Unable to start mocked server, err = %s", err)
	}
}
```

```go
package extapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type (
	// Client describes a client which makes requests to an external API
	Client struct {
		Hostname string
	}

	// HelloResponse describes a successful response to the hello endpoint
	HelloResponse struct {
		Message string `json:"message"`
	}
)

// GetHello calls the hello endpoint from the external API and returns the
// response if successful, or an error otherwise.
func (c Client) GetHello() (*HelloResponse, error) {
	res, err := http.Get(fmt.Sprintf("%s/hello", c.Hostname))
	if err != nil {
		return nil, fmt.Errorf("unable to get extapi hello response, err = %s", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read hello response, err = %s", err)
	}

	decoded := &HelloResponse{}
	if err := json.Unmarshal(body, decoded); err != nil {
		return nil, fmt.Errorf("unable to decode hello response, err = %s", err)
	}

	return decoded, nil
}
```

```go
package extapi_test

import (
	"testing"

	"github.com/alr-lab/ptp/extapi"
)

const want = "Hello, world!"

func TestClient(t *testing.T) {
	// Arrange
	c := &extapi.Client{Hostname: "http://mockapi:8081"}

	// Act
	res, err := c.GetHello()

	// Assert
	if err != nil {
		t.Fatalf("Unable to get mockapi hello response, err = %s", err)
	}
	if res.Message != want {
		t.Fatalf("got %q, want %q", res.Message, want)
	}
}
```

## Contract tests

- Apply to a microservice context
- Consumer/Provider
- Subscriber/Publisher (event-driven communications)
- Specify interfaces between services where consumer and provider are
  spread
- Contract tests ensure implementations on consumer and provider fulfill
  contract
- Act as regression test suite
- Observe deviations early
- _Consumer-Driven Contract tests (CDC tests)_ implement the contract at a
  consumer level
  1. Consuming team writes tests with their expectations
  2. Consuming team shares tests with the providing team
  3. Providing team runs CDC tests continuously
  4. Communications start again when tests are failing
- CDC tests are a step towards establishing autonomous teams

## UI tests

## End-to-end tests

## Acceptance tests

## Exploratory testing

## Terminology

## Deployment pipeline

## Test duplication

## Clean test code
