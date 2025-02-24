package service

import (
	"dc_honest/src/internal/core/ports/output"
)

type ShuffleService struct {
	repo output.ShuffleRepoPort
}

func NewShuffleService(repo output.ShuffleRepoPort) *ShuffleService {
	return &ShuffleService{
		repo: repo,
	}
}

func (s *ShuffleService) ShuffleDeck(id string, clientID string) error {
	return s.repo.ShuffleDeck(id, clientID)
}

func (s *ShuffleService) ShuffleLevel(id string, clientID string) error {
	return s.repo.ShuffleLevel(id, clientID)
}
