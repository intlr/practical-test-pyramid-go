// +build ui

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
