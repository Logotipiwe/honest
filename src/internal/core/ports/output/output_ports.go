package output

import (
	"database/sql"
	"dc_honest/src/internal/core/domain"
)

type DecksStoragePort interface {
	GetAvailableDecks(clientID string) ([]domain.Deck, error)
	UnlockDeck(clientID, deckID string) error
}

type ShuffleRepoPort interface {
	ShuffleDeck(id string, clientID string) error
	ShuffleLevel(id string, clientID string) error
}

type QuestionRepoPort interface {
	GetRandomQuestion(tx *sql.Tx, clientID string, levelID string) (domain.Question, error)
	AddQuestionToHistory(tx *sql.Tx, clientID string, question domain.Question) error
	AddQuestionToUsedQuestions(tx *sql.Tx, clientID string, question domain.Question) error
	AreAllQuestionsGot(tx *sql.Tx, clientID string, levelID string) (bool, error)
	ClearUsedQuestions(tx *sql.Tx, clientID, levelID string) error
}

type LevelsRepoPort interface {
	LevelExists(levelID string) (bool, error)
}
