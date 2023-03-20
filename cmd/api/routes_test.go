package main

import (
	"net/http"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_application_routes(t *testing.T) {

	var registered = []struct {
		route  string
		method string
	}{
		{"/createDeck", "POST"},
		{"/openDeck/{deckId}", "GET"},
		{"/drawDeck/{deckId}", "GET"},
	}

	mux := app.routes()

	chiRoutes := mux.(chi.Routes)

	for _, route := range registered {
		//check to see if the route exists
		if !routeExists(route.route, route.method, chiRoutes) {
			t.Errorf("route %s is not registered", route.route)
		}
	}
}

func routeExists(testRoute, testMethod string, chiRoutes chi.Routes) bool {
	found := false

	_ = chi.Walk(chiRoutes, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if strings.EqualFold(method, testMethod) && strings.EqualFold(route, testRoute) {
			found = true
		}
		return nil
	})
	return found
}
