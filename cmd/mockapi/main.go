package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, _ *http.Request) {
		log.Print("Serving request")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		io.WriteString(w, `{"message":"Hello, world!"}`)
	})

	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Unable to start mocked server, err = %s", err)
	}
}
