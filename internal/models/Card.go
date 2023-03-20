package models

import "strconv"

type Card struct {
	ID      string `json:"-"`
	Value   string `json:"value"`
	Suit    string `json:"suit"`
	Code    string `json:"code"`
	DeckId  string `json:"-"`
	IsDrawn bool   `json:"-"`
}

func (c *Card) SetCode() {
	var suitCode = c.Suit[0:1]
	var valueCode = c.Value[0:1]
	var firstCharacter = ""

	if isNumeric(valueCode) {
		firstCharacter = c.Value
	} else {
		firstCharacter = valueCode
	}

	c.Code = string(firstCharacter + suitCode)
}

func isNumeric(s string) bool {
	_, err := strconv.ParseInt(s, 10, 32)
	return err == nil
}
