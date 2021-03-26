// Copyright 2021 Alexandre Le Roy. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in
// the LICENSE file.

// +build unit

/*

Simple unit test

Part of the Test Double repository I also published
https://github.com/alr-lab/test-double-go

The System Under Test is the ``Service'' object, and both the public
interfaces ``New'' and ``Get'' are tested. It is considered a solitary
unit test as we are doubling the datastore.

*/
package service_test

import (
	"context"
	"testing"

	"github.com/alr-lab/practical-test-pyramid-go/pkg/service"
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
