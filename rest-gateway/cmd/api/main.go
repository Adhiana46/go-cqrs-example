package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	appName    = "REST-Gateway"
	appVersion = "1.0"
	port       = "80"
)

type Config struct {
	AppName    string
	AppVersion string
}

func main() {
	app := Config{
		AppName:    appName,
		AppVersion: appVersion,
	}

	log.Printf("Starting %s service on port %s\n", appName, port)

	s := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	// starting the server
	if err := s.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
