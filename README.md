# Practical Test Pyramid with Go

A [Github Page][github-page] is available

## Getting Started

```txt
$ docker-compose up -d --build
```

## Test

```txt
$ docker exec -it app go test -v -p=1 ./...
```

```txt
?   	github.com/alr-lab/practical-test-pyramid-go/cmd/mockapi	[no test files]
?   	github.com/alr-lab/practical-test-pyramid-go/cmd/server	[no test files]
=== RUN   TestClient
--- PASS: TestClient (0.00s)
PASS
ok  	github.com/alr-lab/practical-test-pyramid-go/extapi	(cached)
?   	github.com/alr-lab/practical-test-pyramid-go/internal/dbtesting	[no test files]
=== RUN   TestService
--- PASS: TestService (0.00s)
=== RUN   TestService_fixtures
=== RUN   TestService_fixtures/Successful
=== RUN   TestService_fixtures/Unexisting
--- PASS: TestService_fixtures (0.24s)
    --- PASS: TestService_fixtures/Successful (0.09s)
    --- PASS: TestService_fixtures/Unexisting (0.15s)
PASS
ok  	github.com/alr-lab/practical-test-pyramid-go/service	0.240s
=== RUN   TestStore
=== RUN   TestStore/Valid_customer_identifier_will_return_valid_email
=== RUN   TestStore/Invalid_customer_identifier_will_return_empty_string
--- PASS: TestStore (0.08s)
    --- PASS: TestStore/Valid_customer_identifier_will_return_valid_email (0.00s)
    --- PASS: TestStore/Invalid_customer_identifier_will_return_empty_string (0.00s)
PASS
ok  	github.com/alr-lab/practical-test-pyramid-go/store	0.082s
```

[github-page]: https://alr-lab.github.io/practical-test-pyramid-go
