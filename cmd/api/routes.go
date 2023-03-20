package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {

	mux := chi.NewRouter()

	//register middleware
	mux.Use(middleware.Recoverer)

	//register routes
	mux.Post("/createDeck", app.CreateDeck)
	mux.Get("/openDeck/{deckId}", app.OpenDeck)
	mux.Get("/drawDeck/{deckId}", app.DrawDeck)

	return mux
}
