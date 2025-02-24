package output

import "dc_honest/src/internal/core/domain"

type DecksStoragePort interface {
	GetDecksForClient(clientID string) ([]domain.Deck, error)
}

type ShuffleRepoPort interface {
	ShuffleDeck(id string, clientID string) error
	ShuffleLevel(id string, clientID string) error
}
