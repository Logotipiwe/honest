package input

import "dc_honest/src/internal/core/domain"

type DecksPort interface {
	GetDecksForMainPage(clientID string) ([]domain.Deck, error)
}

type ShufflePort interface {
	ShuffleDeck(id string, clientID string) error
	ShuffleLevel(id string, clientID string) error
}
