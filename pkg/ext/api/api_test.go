// Copyright 2021 Alexandre Le Roy. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in
// the LICENSE file.

// +build integration

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
