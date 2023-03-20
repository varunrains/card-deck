package models

type Deck struct {
	ID        string `json:"deck_id"`
	Shuffled  bool   `json:"shuffled"`
	Remaining int    `json:"remaining"`
	Cards     []Card `json:"cards,omitempty"`
}

func (d *Deck) SetRemaining() {
	d.Remaining = len(d.Cards)
}

func (d *Deck) SetCards(cards []Card) {
	d.Cards = cards
}
