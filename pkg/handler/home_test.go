// Copyright 2021 Alexandre Le Roy. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in
// the LICENSE file.

// +build unit

package handler_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alr-lab/practical-test-pyramid-go/pkg/ext/api"
	"github.com/alr-lab/practical-test-pyramid-go/pkg/handler"
)

type (
	// describes a client mock which will return a successful response
	successfulAPIMock struct{}

	// describes a client mock which will return an erroring response
	erroringAPIMock struct{}
)

func Test_Home(t *testing.T) {
	tt := map[string]struct {
		mock interface {
			GetHello() (*api.HelloResponse, error)
		}
		want   string
		status int
	}{
		"Successful request": {
			mock:   successfulAPIMock{},
			want:   `{"message": "foo"}`,
			status: http.StatusOK,
		},
		"Erroring request": {
			mock:   erroringAPIMock{},
			want:   `{"error": "unable to get message"}`,
			status: http.StatusInternalServerError,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatalf("Unable to create home request, err = %s", err)
			}

			rec := httptest.NewRecorder()
			f := handler.Home(tc.mock)
			f(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tc.status {
				t.Errorf("invalid status, got %d, want %d", res.StatusCode, tc.status)
			}

			raw, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("Unable to read home response, err = %s", err)
			}

			if string(raw) != tc.want {
				t.Errorf("invalid response body, got %q, want %q", raw, tc.want)
			}
		})
	}
}

func (m successfulAPIMock) GetHello() (*api.HelloResponse, error) {
	return &api.HelloResponse{Message: "foo"}, nil
}

func (m erroringAPIMock) GetHello() (*api.HelloResponse, error) {
	return nil, fmt.Errorf("some error")
}
