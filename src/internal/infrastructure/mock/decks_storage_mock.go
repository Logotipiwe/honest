package mock

import (
	"dc_honest/src/internal/core/domain"
)

type DecksStorageMock struct {
	decks                  []domain.Deck
	GetAvailableDecksCalls [][]any
}

func NewDecksStorageMock() *DecksStorageMock {
	return &DecksStorageMock{
		decks:                  make([]domain.Deck, 0),
		GetAvailableDecksCalls: make([][]any, 0),
	}
}

func (d *DecksStorageMock) GetAvailableDecks(clientID string) ([]domain.Deck, error) {
	d.GetAvailableDecksCalls = append(d.GetAvailableDecksCalls, []any{clientID})
	return d.decks, nil
}

func (d *DecksStorageMock) UnlockDeck(clientID, deckID string) error {
	return nil
}

func (d *DecksStorageMock) SetDecks(decks []domain.Deck) {
	d.decks = decks
}

func (d *DecksStorageMock) Clean() {
	d.decks = make([]domain.Deck, 0)
	d.GetAvailableDecksCalls = make([][]any, 0)
}
