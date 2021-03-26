package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alr-lab/practical-test-pyramid-go/pkg/ext/api"
	"github.com/alr-lab/practical-test-pyramid-go/pkg/handler"
)

const (
	// describes home endpoint
	endpointHome = "/"

	// describes the environment variable holding the application port
	envApplicationPort = "APPLICATION_PORT"

	// describes the environment variable holding the external API hostnam
	envExternalAPIHostname = "EXTERNAL_API_HOST"
)

var (
	// describes the application port
	applicationPort = os.Getenv(envApplicationPort)

	// describes the external API hostname
	externalAPIHostname = os.Getenv(envExternalAPIHostname)
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

	s.HandleFunc(endpointHome, handler.HomeHandler(c))

	return s
}
