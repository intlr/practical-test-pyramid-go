---
permalink: /
---
> Personal notes on [the article Ham Vocke wrote][ptp-article] on the
> subject and interpretation to a Golang application &mdash; Alexandre Le
> Roy

The [Practical Test Pyramid with Go][ptp-go] repository hosts all the
examples put together to describe a microservice application. It also
provides the source code for this Github Page.

## Table of Contents

- [Unit tests](#unit-tests)
- [Integration tests](#integration-tests)
  - [Database integration](#database-integration)
  - [External API integration](#external-api-integration)
- [Contract tests](#contract-tests)
- [UI tests](#ui-tests)
- [End-to-end tests](#end-to-end-tests)
- [Acceptance tests](#acceptance-tests)
- [Exploratory testing](#exploratory-testing)
- [Deployment pipeline](#deployment-pipeline)
- [Test duplication](#test-duplication)
- [Clean test code](#clean-test-code)

## Unit tests

Unit tests consist of the foundation of the test suite. It ensures the
_System Under Test (SUT)_ works as intended. Unit tests are fast, therefor
they should prevail other types of tests.

There are two types of unit tests: the _Solitary unit tests_, and the
_Sociable unit tests_. The former describe tests doubling all
collaborators, while the latter describe tests allowing communications with
the real collaborators.

The _Test-Driven Development (TDD)_ lets unit tests guide the development.

Ham Vocke also provides rules of thumb for writing unit tests.

- Write one test class per production class
- Unit test at least public interfaces
- Include happy cases and edge cases, without being too tied to implementation
- Arrange, Act, Assert

### Implementation

Let's dig into a simple implementation we will unit test. I decided to
implement an _application service_ which relies on a datastore to fetch
customer data.

Notice we do not integrate the datastore. At this level of the pyramid
there is no need to be specific. The decision to use a specific datastore
can be done later. Therefor we are only providing a contract for future
implementations.

```go
// The ``service'' package provides a basic service with a dependency on a
// datastore. The application uses this service to handle requests and
// return a simple response.
package service

import "context"

type (
	// Service describes the application service
	Service struct {
		store Store
	}

	// Store defines a contract for a datastore
	Store interface {
		// GetCustomerEmail returns a customer email address from a
		// customer identifier
		GetCustomerEmail(ctx context.Context, id int) string
	}
)

// New returns a new application service
func New(store Store) Service {
	return Service{store: store}
}

// Get returns a specific customer email
func (s Service) Get(ctx context.Context) string {
	return s.store.GetCustomerEmail(ctx, 42)
}
```

### Test

Because we didn't integrated the datastore yet, we used a contract which is
allowing us to inject any kind of datastore to the service.

When testing the service, this allows us to actually provide a stub. At
this level of test we don't care about testing the integration with the
datastore. We only care that the service is working as intended, that it
returns the correct value.

```go
package service_test

import (
	"context"
	"testing"

	"github.com/alr-lab/test-double-go/pkg/service"
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

Using external tests we force ourselves to write tests only for public
functions. This is done by using another package than the actual
implementation - `service_test` instead of `service`.

We provide a stub to the service as the datastore and we make sure the
service returns the valid data. Because we doubled all dependencies of the
application service, this unit test can be considered as a _Solitary unit
test_.

This simple implementation only have a happy path. This is because the
datastore actually handles errors and will return an empty string in case
something goes bad. This is not what one may expect for a production
scenario but it allows us to focus on experimenting with unit tests.

Notice the _Arrange, Act, Assert_ pattern. This will be repeated all over
the pyramid.

I cover the types of double in Golang with the [Test Double with
Go][test-double-go] repository

## Integration tests

Integration tests ensure integrations with other parts of the system are
working. Those parts can be databases, filesystems, networks... In order to
do so, we may need to run those components we are integrating.

There are two types of integration tests: one where we test through the
entire stack, and one we test the integrations one by one, doubling the
others if needed.

Because we test the application with the actual integration and not a
double, those tests are slower compared to the unit tests. That is why
those tests should be less common.

### Database integration

Most of database integration tests consist of the following steps.

1. Start the database
2. Connect the application to the database
3. Interact with the database
4. Validate the expectations

Before to get into the implementation of such an integration, let's have a
look at the test helper which allows us to create database connections.
This will be extended later to include fixtures, but for now, we need a way
to simply communicate with a datastore.

```go
package dbtesting

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql" // required to use the MySQL driver
)

// DatabaseHelper is a test helper. It doesn't start with the `Test`
// prefix, thus will not be executed by the Go runtime when testing the
// application. However, we will call this function from a test function,
// passing the test argument to be able to fail the test at any point.
func DatabaseHelper(t *testing.T) *sql.DB {
	t.Helper()

	// Open MySQL connection and fill database with application schema...
	// https://github.com/alr-lab/practical-test-pyramid-go/blob/master/internal/dbtesting/mysql.go

	return conn
}
```

Now let's dive into the implementation of the database integration. We need
an object on which we can query the database and get the data we need.

```go
package store

import (
	"context"
	"database/sql"
	"fmt"
)

// Store describes a datastore
type Store struct{
	conn *sql.DB
}

// SetConn sets the database connection to the datastore object
func (st *Store) SetConn(conn *sql.DB) *Store {
	st.conn = conn
	return st
}

// GetCustomerEmail returns the customer email address.
//
// Notice this method was extended to also accept a context as one
// parameter, and to return an error.
func (st *Store) GetCustomerEmail(ctx context.Context, id int) (string, error) {
	query := `SELECT email FROM Customers WHERE id = ?`

	args := []interface{}{id}
	res, err := st.conn.QueryContext(ctx, query, args...)
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

Using `docker-compose`, we [spin up a MySQL database][docker-compose-file].
We are now able to test the integration.

```go
package store_test

import (
	"context"
	"testing"

	"github.com/alr-lab/practical-test-pyramid-go/internal/dbtesting"
	"github.com/alr-lab/practical-test-pyramid-go/pkg/ext/store"
)

func TestStore(t *testing.T) {
	// Arrange
	conn := dbtesting.DatabaseHelper(t)
	defer func() { _ = conn.Close() }()
	st := (&store.Store{}).SetConn(conn)
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

Notice the _Arrange, Act, Assert_ pattern? The _Arrange_ step is bigger but
it prepares more edge cases than unit-tests. By doing so, we ensure our
integration will fail properly when needed.

### External API integration

Similar to database integration tests, most external API integration tests
consist of specific tests.

1. Start the application
2. Start the instance of the API
3. Interact with the API
4. Validate the expectations

As we don't have an instance of the external API, we are creating a mock,
returning the response we are expecting.

```go
// The ``mockapi'' application is a simple server acting as a replacement
// for an external API.
//
// With a ``docker-compose'' file we can start the server and allow our
// application to test its production-integration.
package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	log.Print("Starting server")
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

Let's look at the implementation of the external API integration. We need
a mechanism to call the external API, request a specific endpoint, receive
a response, parse this response, and output something from it.

```go
package api

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
		return nil, fmt.Errorf("unable to get external API hello response, err = %s", err)
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

Now, testing is easier. We create a basic client that will hit the mocked
API and test the integration as we would test a database integration.

```go
package api_test

import (
	"testing"

	"github.com/alr-lab/practical-test-pyramid-go/pkg/ext/api"
)

const want = "Hello, world!"

func TestClient(t *testing.T) {
	// Arrange
	c := &api.Client{Hostname: "http://mockapi:8081"}

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

Contract tests applies to a microservice context, where two microservices
are communicating with each other. In that scenario, there is a _Consumer_
and a _Provider_. In case of event-driven communications, we call the
former _Subscriber_ and the latter _Publisher_.

The contract tests ensure both implementations fulfill a contract
previously defined. It acts as a regression test suite and allow to observe
deviations early.

When this contract is provided by the consumer, we call the tests
_Consumer-Driven Contract tests (CDC tests)_. Such a mechanism is
implemented by following those steps.

1. The consuming team writes tests with its expectations
2. The consuming team shares tests with the providing team
3. The providing team runs CDC tests continuously
4. Communications start again when tests are failing

One popular way to implement contract testing with Go would be to adopt
Pact in our pipeline. This is a very deep topic and will be covered in
another article.

## UI tests

UI tests are more than testing web browser interfaces. Think of REST API,
CLI... they all have interfaces.

UI tests are about testing the user interface is working as expected. But
testing UI can be done in a modular way, as testing JavaScript code with
the backend being stubbed.

Think of testing the behaviour, the layout, the usability...

In our case, it is as simple as calling the API that we started using
docker-compose, and expecting the response to implement a specific format.
There is no point to expect the response to be precise, simply to fulfill
some kind of format we define during the UI testing.

```go
package main_test

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

type response struct {
	Message string `json:"message"`
}

func Test(t *testing.T) {
	// Arrange
	req, err := http.NewRequest("GET", "http://app:8080/", nil)
	if err != nil {
		t.Fatalf("Unable to create application request, err = %s", err)
	}

	// Act
	res, err := (http.DefaultClient).Do(req)

	// Assert
	if err != nil {
		t.Fatalf("Unable to request application, err = %s", err)
	}
	defer res.Body.Close()

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Unable to read response, err = %s", err)
	}

	var decoded response
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("Unable to decode response, err = %s", err)
	}
}
```

## End-to-end tests

End-to-end tests are about testing the fully-integrated system.

Those tests are heavy, they may fail for unexpected reasons such as
timeouts... They require maintenance and run slowly. This is why one would
aim to reduce end-to-end tests to the minimum.

Let's dig into implementing an end-to-end test. First, we need a way to get
a clean version of the database with a specific state of data. This is why
the DatabaseHelper was updated to accept a fixture directory.

```go
package dbtesting

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/testfixtures.v2"
)

func DatabaseHelper(t *testing.T, fixtureDir string) *sql.DB {
	t.Helper()

	// Open MySQL connection and fill database with testing data...

	if fixtureDir != "" {
		testfixtures.SkipDatabaseNameCheck(true)
		fixtures, err := testfixtures.NewFolder(conn, &testfixtures.MySQL{}, fixtureDir)
		if err != nil {
			t.Fatalf("Unable to create fixtures from fixture directory, err = %s", err)
		}

		if err := fixtures.Load(); err != nil {
			t.Fatalf("Unable to load fixtures, err = %s", err)
		}
	}

	return conn
}
```

Then, we are able to define Yaml fixtures. As an example, `Customers.yaml`
will allow the DatabaseHelper to request `testfixtures` to fill the
Customers table with the following data.

```yaml
- id: 42
  email: foo
```

Then, we test the service, we could go one step further and call the API as
we did, but it gives us less flexibility in controlling how to setup the
application.

For each test, we start a fresh database version, setup the service to use
the database instance, and call the service. We then assert the response is
what we would expect for the given scenario.

```go
package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/alr-lab/practical-test-pyramid-go/internal/dbtesting"
	"github.com/alr-lab/practical-test-pyramid-go/pkg/ext/store"
	"github.com/alr-lab/practical-test-pyramid-go/pkg/service"
)

func TestService_fixtures(t *testing.T) {
	// Arrange
	tt := map[string]struct {
		want     string
		fixtures string
		ctx      context.Context
	}{
		"Successful": {
			want:     "foo",
			fixtures: "successful",
			ctx:      context.Background(),
		},
		"Unexisting": {
			want:     "",
			fixtures: "unexisting",
			ctx:      context.Background(),
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			conn := dbtesting.DatabaseHelper(t, fmt.Sprintf("fixtures/%s", tc.fixtures))
			defer conn.Close()
			st := (&store.Store{}).SetConn(conn)
			serv := service.New(st)

			// Act
			got := serv.Get(tc.ctx)

			// Assert
			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}
```

## Acceptance tests

Acceptance tests are tests ensuring the application works from a user's
perspective. Think of the _Given, When, Then_ mnemonic.

* Given a user is connected
* When a user go to the homepage
* Then the page should display the username

The _Behaviour-Driven Development (BDD)_ helps to focus on this user's
perspective.

There are different levels of granularity to describe acceptance tests. You
may test the user's perspective through the user interface, but you can
also test a feature works properly before to reach the user interface.

From a user's perspective, we already proved that our application is able
to handle requests and respond the way we want. Both the UI test and
end-to-end test proved that. We can confidently skip this step.

## Exploratory testing

Exploratory testing is everything related to manual testing an application.
Ham Vocked wrote in his article that the best way to reach confidence with
those tests is to adopt a destructive mindset, to try to break the
application.

Documenting while testing is a good way to keep records of the things you
may discover during the process while not loosing focus.

## Deployment pipeline

- Automated pipeline using _Continuous Integration (CI)_ or _Continuous
  Delivery (CD)_ will provide gradual confidence the application is ready
  to be deployed to production
- A good pipeline breaks as early as possible
- Fast tests early, slow tests later

## Test duplication

- Avoid test duplication throughout the different layers of the pyramid
- Writing and maintaining tests takes time
- Reading and running tests also takes time
- Rules of thumb
  1. If higher-level test detects an errors, a lower-level test is needed
  2. Push tests as far down as possible

## Clean test code

1. Test code is as important as production code
2. Test one condition per test
3. Arrange, Act, Assert

[docker-compose-file]: https://github.com/alr-lab/practical-test-pyramid-go/blob/master/docker-compose.yml
[ptp-article]: https://martinfowler.com/articles/practical-test-pyramid.html
[ptp-go]: https://github.com/alr-lab/practical-test-pyramid-go
[test-double-go]: https://github.com/alr-lab/test-double-go
