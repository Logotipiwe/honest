package ms

import "database/sql"

type ShuffleRepoMs struct {
	db *sql.DB
}

func NewShuffleRepoMs(db *sql.DB) *ShuffleRepoMs {
	return &ShuffleRepoMs{
		db: db,
	}
}

// TODO покрыть тестами

func (s *ShuffleRepoMs) ShuffleDeck(id string, clientID string) error {
	sql := `DELETE FROM used_questions
		WHERE question_id IN (
		    SELECT q.id from questions q 
		    LEFT JOIN levels l on q.level_id = l.id
		    WHERE l.deck_id = ?
		)
		AND client_id = ?
	`
	_, err := s.db.Exec(sql, id, clientID)
	return err
}

func (s *ShuffleRepoMs) ShuffleLevel(id string, clientID string) error {
	sql := `DELETE FROM used_questions
		WHERE question_id IN (
		    SELECT q.id from questions q 
		    WHERE q.level_id = ?
		)
		AND client_id = ?
	`
	_, err := s.db.Exec(sql, id, clientID)
	return err
}
