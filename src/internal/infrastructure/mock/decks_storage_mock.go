package mock

import (
	"dc_honest/src/internal/core/domain"
)

type DecksStorageMock struct {
	decks []domain.Deck
}

func NewDecksStorageMock() *DecksStorageMock {
	return &DecksStorageMock{
		decks: make([]domain.Deck, 0),
	}
}

func (d *DecksStorageMock) GetDecksForClient(clientID string) ([]domain.Deck, error) {
	return d.decks, nil
}

func (d *DecksStorageMock) SetDecks(decks []domain.Deck) {
	d.decks = decks
}

func (d *DecksStorageMock) Clean() {
	d.decks = make([]domain.Deck, 0)
}
