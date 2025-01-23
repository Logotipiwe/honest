package output

import "dc_honest/src/internal/core/domain"

type DecksStoragePort interface {
	GetDecksForClient(clientID string) ([]domain.Deck, error)
}
