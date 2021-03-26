package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alr-lab/practical-test-pyramid-go/pkg/ext/api"
	"github.com/alr-lab/practical-test-pyramid-go/pkg/handler"
)

var (
	// describes the application port
	applicationPort = os.Getenv("APPLICATION_PORT")

	// describes the external API hostname
	externalAPIHostname = os.Getenv("EXTERNAL_API_HOST")
)

func main() {
	c := &api.Client{Hostname: externalAPIHostname}
	addr := fmt.Sprintf(":%s", applicationPort)

	err := http.ListenAndServe(addr, handle(c))
	if err != nil {
		log.Fatalf("Unable to start server, err = %s", err)
	}
}

func handle(c *api.Client) http.Handler {
	s := http.NewServeMux()

	s.HandleFunc("/", handler.Home(c))

	return s
}
