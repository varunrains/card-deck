package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/varunrains/carddeck/internal/models"
)

const SHUFFLED string = "shuffled"
const CARDS string = "cards"
const COUNT string = "count"

func (app *application) CreateDeck(w http.ResponseWriter, r *http.Request) {

	// Parse query parameters
	shuffled := false
	cards := []string{}
	queryParams := r.URL.Query()

	if _, ok := queryParams[SHUFFLED]; ok && queryParams.Get(SHUFFLED) == "true" {
		shuffled = true
	}

	if _, ok := queryParams[SHUFFLED]; ok {
		_, err := strconv.ParseBool(queryParams.Get(SHUFFLED))
		if err != nil {
			app.errorJSON(w, errors.New("please specify the proper value for shuffled"+err.Error()), http.StatusBadRequest)
			return
		}
	}

	if _, ok := queryParams[CARDS]; ok {
		cards = strings.Split(queryParams.Get(CARDS), ",")
	}

	deck, err := app.DB.CreateDeck(shuffled, cards)

	if err != nil {
		app.errorJSON(w, errors.New("error occured while creating the deck"+err.Error()), http.StatusInternalServerError)
		return
	}

	deck.SetRemaining()
	deck.SetCards(nil)
	app.writeJSON(w, http.StatusOK, deck)

}

func (app *application) OpenDeck(w http.ResponseWriter, r *http.Request) {

	deckId := chi.URLParam(r, "deckId")

	deck, err := app.DB.OpenDeck(deckId)
	if err != nil {
		app.errorJSON(w, errors.New("error occured while opening the deck :"+err.Error()), http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, deck)
}

func (app *application) DrawDeck(w http.ResponseWriter, r *http.Request) {

	deckId := chi.URLParam(r, "deckId")
	queryParams := r.URL.Query()
	n := queryParams.Get(COUNT)

	if len(n) == 0 {
		app.errorJSON(w, errors.New("specify the number of cards to draw"), http.StatusBadRequest)
		return
	}

	nc, err := strconv.ParseInt(n, 10, 8)

	if err != nil {
		app.errorJSON(w, errors.New("specify proper number of cards to draw"), http.StatusBadRequest)
		return
	}

	deckCards, err := app.DB.DrawDeck(deckId, int(nc))
	if err != nil {
		app.errorJSON(w, errors.New("error occured while drawing the deck : "+err.Error()), http.StatusInternalServerError)
		return
	}
	cards := []models.Card{}
	cards = append(cards, *deckCards...)
	app.writeJSON(w, http.StatusOK, cards)
}
