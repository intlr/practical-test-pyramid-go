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

	"github.com/alr-lab/ptp/service"
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
