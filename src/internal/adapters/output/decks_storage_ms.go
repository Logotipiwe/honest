package output

import (
	"database/sql"
	"dc_honest/src/internal/core/domain"
	"dc_honest/src/internal/infrastructure/ms"
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
		var deck ms.DeckModel
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

func NewDecksStorageMs(
	db *sql.DB,
) *DecksStorageMs {
	return &DecksStorageMs{
		db: db,
	}
}
