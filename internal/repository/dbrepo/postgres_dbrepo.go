package dbrepo

import (
	"context"

	"database/sql"
	"errors"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/varunrains/carddeck/internal/models"
)

type PostgresDBRepo struct {
	DB *sql.DB
}

const dbTimeout = time.Second * 5

var values = []string{"ACE", "2", "3", "4", "5", "6", "7", "8", "9", "10", "JACK", "QUEEN", "KING"}
var suits = []string{"SPADES", "DIAMONDS", "CLUBS", "HEARTS"}

func (m *PostgresDBRepo) Connection() *sql.DB {
	return m.DB
}

func (m *PostgresDBRepo) CreateDeck(shuffled bool, cards []string) (*models.Deck, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var deckId = uuid.NewString()
	// start a transaction
	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	sqlInsert := "INSERT INTO deck (id, shuffled) VALUES ($1, $2)"

	_, err = tx.ExecContext(ctx, sqlInsert, deckId, shuffled)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	deck := createDeckModel(shuffled, cards, deckId)

	sqlInsert = "INSERT INTO card (id, deck_id, value, suit) VALUES ($1, $2, $3, $4)"

	stmtInsert, err := tx.PrepareContext(ctx, sqlInsert)
	if err != nil {
		return nil, err
	}
	defer stmtInsert.Close()

	for _, card := range deck.Cards {
		_, err = stmtInsert.ExecContext(ctx, uuid.New().String(), deckId, card.Value, card.Suit)
		if err != nil {
			return nil, err
		}
	}

	tx.Commit()

	return &deck, nil

}

func (m *PostgresDBRepo) OpenDeck(deckId string) (*models.Deck, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		select 
				d.id, d.shuffled, c.id, c.value, c.suit
			from deck d
		join 
			card 
				c on d.id = c.deck_id
		where 
			d.id= $1 and (c.isDrawn is null or c.isDrawn = $2)`

	rows, err := m.DB.QueryContext(ctx, query, deckId, "false")

	if err != nil {
		return nil, err
	}

	cards := []models.Card{}
	deck := models.Deck{}

	for rows.Next() {
		var c models.Card

		err := rows.Scan(
			&deck.ID,
			&deck.Shuffled,
			&c.ID,
			&c.Value,
			&c.Suit,
		)
		c.SetCode()
		if err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}
	defer rows.Close()

	if deck.ID == uuid.Nil.String() || deck.ID == "" {
		return nil, errors.New("not a valid deck id/ deck is empty")
	}

	deck.SetCards(cards)
	deck.SetRemaining()

	return &deck, nil

}

func (m *PostgresDBRepo) DrawDeck(deckId string, count int) (*[]models.Card, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var remaining int
	query := `
		select 
				d.id, d.shuffled, c.id, c.value, c.suit
			from deck d
		join 
			card 
				c on d.id = c.deck_id
		where 
			d.id= $1
		LIMIT $2`

	countQuery := `select count(*) from card where deck_id = $1 and (isDrawn is null or isDrawn = $2) group by deck_id`
	// start a transaction
	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	nr := tx.QueryRowContext(ctx, countQuery, deckId, "false")
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = nr.Scan(&remaining)

	if err == sql.ErrNoRows {
		return nil, errors.New("not a valid deck id")
	}

	if err != nil || remaining < count {
		return nil, errors.New("remaining cards in the deck is less than the count")
	}

	rows, err := tx.QueryContext(ctx, query, deckId, count)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	cards := []models.Card{}
	deck := models.Deck{}

	for rows.Next() {
		var c models.Card

		err := rows.Scan(
			&deck.ID,
			&deck.Shuffled,
			&c.ID,
			&c.Value,
			&c.Suit,
		)
		c.SetCode()
		if err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}
	deck.SetCards(cards)
	defer rows.Close()

	sqlInsert := "UPDATE CARD SET isDrawn = $1 WHERE ID = $2"

	stmtInsert, err := tx.PrepareContext(ctx, sqlInsert)
	if err != nil {

		return nil, err
	}
	defer stmtInsert.Close()

	for _, card := range deck.Cards {
		_, err = stmtInsert.ExecContext(ctx, "true", card.ID)
		if err != nil {

			return nil, err
		}
	}

	tx.Commit()

	return &deck.Cards, nil
}

func shuffle(shuffled bool, deck *models.Deck) {
	if shuffled {
		rand.Shuffle(len(deck.Cards), func(i, j int) {
			deck.Cards[i], deck.Cards[j] = deck.Cards[j], deck.Cards[i]
		})
	}
}

func contains(cards []string, card models.Card) bool {
	card.SetCode()
	for _, c := range cards {
		if c == card.Code {
			return true
		}
	}
	return false
}

func createDeckModel(shuffled bool, cards []string, deckId string) models.Deck {

	deck := models.Deck{
		ID:       deckId,
		Cards:    make([]models.Card, 0, 52),
		Shuffled: shuffled,
	}

	for _, suit := range suits {
		for _, value := range values {
			card := models.Card{
				Value: value,
				Suit:  suit,
			}

			if len(cards) == 0 || contains(cards, card) {
				deck.Cards = append(deck.Cards, card)
			}
		}
	}

	shuffle(shuffled, &deck)

	return deck
}
