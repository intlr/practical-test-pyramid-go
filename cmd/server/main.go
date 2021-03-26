// Copyright 2021 Alexandre Le Roy. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in
// the LICENSE file.

/*

The server application acts as a server handling customer requests. The only
available endpoint so far is the home endpoint, which greats customers. The
server uses a MySQL data store to access customer emails and fetches a
greating message from an external API.

This application allows us to explore testing at the different layers of the
practical test pyramid introduced by Ham Vocke.

*/
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
	if applicationPort == "" {
		log.Fatal("Missing application port")
	}

	if externalAPIHostname == "" {
		log.Fatal("Missing external API hostname")
	}

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
