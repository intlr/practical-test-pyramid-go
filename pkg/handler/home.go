package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/alr-lab/practical-test-pyramid-go/pkg/ext/api"
)

// describes the client contract
type client interface {
	GetHello() (*api.HelloResponse, error)
}

// HomeHandler returns the handler function responsible to handle home
// requests
func HomeHandler(c client) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		log.Print("Serving root request")

		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		res, err := c.GetHello()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, `{"error": "unable to get message"}`)
			return
		}

		io.WriteString(w, fmt.Sprintf(`{"message": "%s"}`, res.Message))
	}
}