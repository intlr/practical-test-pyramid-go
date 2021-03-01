package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/alr-lab/practical-test-pyramid-go/pkg/ext/api"
)

func main() {
	c := (&api.Client{Hostname: "http://mockapi:8081"})
	log.Print("Starting server...")

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		log.Print("Serving request")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		res, err := c.GetHello()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, `{"error": "unable to get message"}`)
			return
		}

		io.WriteString(w, fmt.Sprintf(`{"message": "%s"}`, res.Message))
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Unable to start server, err = %s", err)
	}
}
