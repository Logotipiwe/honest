package output

import (
	"dc_honest/src/internal/core/domain"
)

type DecksStorageMock struct {
	decks map[string][]domain.Deck
}

func NewDecksStorageMock() *DecksStorageMock {
	return &DecksStorageMock{
		decks: make(map[string][]domain.Deck),
	}
}

func (d DecksStorageMock) GetDecksForClient(clientID string) ([]domain.Deck, error) {
	decks, has := d.decks[clientID]
	if !has {
		return []domain.Deck{}, nil
	}
	return decks, nil
}

func (d DecksStorageMock) SetDecksForClient(clientID string, decks []domain.Deck) {
	d.decks[clientID] = decks
}
