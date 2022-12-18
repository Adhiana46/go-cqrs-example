package main

import (
	"fmt"
	"net/http"
)

func (app *Config) WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Welcome to %s version %s", app.AppName, app.AppVersion),
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}
