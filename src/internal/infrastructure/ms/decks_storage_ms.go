package ms

import (
	"database/sql"
	"dc_honest/src/internal/core/domain"
	"fmt"
	"strings"
)

type DecksStorageMs struct {
	db *sql.DB
}

const all = "id, name, description, labels, vector_image_id, hidden, promo"

func (d *DecksStorageMs) GetDecksForClient(clientID string) ([]domain.Deck, error) {
	decks := make([]domain.Deck, 0)
	rows, err := d.db.Query("SELECT " + all + " FROM decks where !hidden")
	if err != nil {
		return decks, err
	}
	defer rows.Close()
	for rows.Next() {
		var deck DeckModel
		_ = rows.Scan(&deck.ID, &deck.Name, &deck.Description, &deck.Labels,
			&deck.Image, &deck.IsHidden, &deck.PromoCode)
		decks = append(decks, deck.ToDeck())
	}
	err = rows.Err()
	if err != nil {
		return decks, err
	}
	return decks, nil
}

func (d *DecksStorageMs) SaveDecks(decks []domain.Deck) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	q := make([]string, 0)
	for range strings.Split(all, ", ") {
		q = append(q, "?")
	}
	for _, deck := range decks {
		_, err := tx.Exec(fmt.Sprintf(
			"INSERT INTO decks (%v) VALUES (%v)",
			all, strings.Join(q, ","),
		), deck.ID, deck.Name, deck.Description, strings.Join(deck.Labels, ";"), deck.Image, deck.IsHidden, deck.PromoCode)
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

func NewDecksStorageMs(
	db *sql.DB,
) *DecksStorageMs {
	return &DecksStorageMs{
		db: db,
	}
}
