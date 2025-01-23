package service

import (
	"dc_honest/src/internal/core/domain"
	"dc_honest/src/internal/core/ports/output"
)

type DecksService struct {
	storage output.DecksStoragePort
}

func NewDecksService(
	storage output.DecksStoragePort,
) *DecksService {
	return &DecksService{
		storage: storage,
	}
}

func (d *DecksService) GetDecksForMainPage(clientID string) ([]domain.Deck, error) {
	decks, err := d.storage.GetDecksForClient(clientID)
	return decks, err
}