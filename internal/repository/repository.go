package repository

import (
	"database/sql"

	"github.com/varunrains/carddeck/internal/models"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	CreateDeck(shuffled bool, cards []string) (*models.Deck, error)
	OpenDeck(deckId string) (*models.Deck, error)
	DrawDeck(deckId string, count int) (*[]models.Card, error)
}
