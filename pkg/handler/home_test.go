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
	// describes the client contract
	client interface {
		GetHello() (*api.HelloResponse, error)
	}

	successfulAPIMock struct{}

	erroringAPIMock struct{}
)

func (m successfulAPIMock) GetHello() (*api.HelloResponse, error) {
	return &api.HelloResponse{Message: "foo"}, nil
}

func (m erroringAPIMock) GetHello() (*api.HelloResponse, error) {
	return nil, fmt.Errorf("some error")
}

func Test_HomeHandler(t *testing.T) {
	tt := map[string]struct {
		mock client
		want string
	}{
		"Successful request": {
			mock: successfulAPIMock{},
			want: `{"message": "foo"}`,
		},
		"Erroring request": {
			mock: erroringAPIMock{},
			want: `{"error": "unable to get message"}`,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatalf("Unable to create home request, err = %s", err)
			}

			rec := httptest.NewRecorder()
			f := handler.HomeHandler(tc.mock)
			f(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			raw, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("Unable to read home response, err = %s", err)
			}

			if string(raw) != tc.want {
				t.Fatalf("got %q, want %q", raw, tc.want)
			}
		})
	}
}
