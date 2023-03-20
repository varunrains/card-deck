package dbrepo

import (
	"database/sql"

	"github.com/varunrains/carddeck/internal/models"
)

type TestDBRepo struct{}

func (m *TestDBRepo) Connection() *sql.DB {
	return nil
}

func (m *TestDBRepo) CreateDeck(shuffled bool, cards []string) (*models.Deck, error) {
	var deck = models.Deck{
		ID:        "1",
		Shuffled:  false,
		Remaining: 1,
		Cards:     []models.Card{},
	}
	return &deck, nil
}

func (m *TestDBRepo) OpenDeck(deckId string) (*models.Deck, error) {
	var deck = models.Deck{
		ID:        "1",
		Shuffled:  false,
		Remaining: 1,
		Cards:     []models.Card{},
	}
	return &deck, nil
}

func (m *TestDBRepo) DrawDeck(deckId string, count int) (*[]models.Card, error) {
	var deck = models.Deck{
		ID:        "1",
		Shuffled:  false,
		Remaining: 1,
		Cards:     []models.Card{},
	}
	return &deck.Cards, nil
}
