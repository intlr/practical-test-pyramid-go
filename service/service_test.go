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
	"fmt"
	"testing"

	"github.com/alr-lab/practical-test-pyramid-go/internal/dbtesting"
	"github.com/alr-lab/practical-test-pyramid-go/service"
	"github.com/alr-lab/practical-test-pyramid-go/store"
)

const email = "fake"

type StubStore struct{}

func (s StubStore) GetCustomerEmail(_ context.Context, _ int) (string, error) {
	return email, nil
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
			st := &store.Store{}
			st.SetConn(conn)
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
