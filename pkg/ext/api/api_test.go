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
