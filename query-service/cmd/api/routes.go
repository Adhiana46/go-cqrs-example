package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		payload := jsonResponse{
			Error:   false,
			Message: fmt.Sprintf("Welcome to %s version %s", app.AppName, app.AppVersion),
		}

		_ = app.writeJSON(w, http.StatusOK, payload)
	})

	mux.Route("/articles", func(r chi.Router) {
		r.Get("/", app.GetArticlesHandler)
		r.Get("/{uuid}", app.GetSingleArticleHandler)
	})

	return mux
}
