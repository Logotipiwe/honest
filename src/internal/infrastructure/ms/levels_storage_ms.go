package ms

import (
	"database/sql"
	"dc_honest/src/internal/core/domain"
	"fmt"
)

type LevelsStorageMs struct {
	db *sql.DB
}

func NewLevelsStorageMs(db *sql.DB) *LevelsStorageMs {
	return &LevelsStorageMs{db: db}
}

func (l *LevelsStorageMs) SaveLevel(level domain.Level) error {
	return l.SaveLevels([]domain.Level{level})
}

func (l *LevelsStorageMs) SaveLevels(levels []domain.Level) error {
	tx, err := l.db.Begin()
	if err != nil {
		return err
	}
	for _, level := range levels {
		_, err := tx.Exec(fmt.Sprintf(
			"INSERT INTO levels (%v) VALUES (%v)",
			"id, deck_id, level_order, name, emoji, color_start, color_end, color_button, description",
			"?,?,?,?,?,?,?,?,?",
		),
			level.ID, level.Deck.ID, level.Order, level.Name, "", "", "", "", level.Description,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
