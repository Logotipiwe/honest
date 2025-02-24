package mock

type ShuffleMock struct {
}

func NewShuffleMock() *ShuffleMock {
	return &ShuffleMock{}
}

func (s *ShuffleMock) ShuffleDeck(id string, clientID string) error {
	return nil
}

func (s *ShuffleMock) ShuffleLevel(id string, clientID string) error {
	return nil
}
