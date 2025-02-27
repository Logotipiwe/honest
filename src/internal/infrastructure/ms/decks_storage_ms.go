package ms

import (
	"database/sql"
	"dc_honest/src/internal/core/domain"
	"fmt"
	"strings"
)

type DecksMsStorage struct {
	db *sql.DB
}

func (d *DecksMsStorage) GetAvailableDecks(clientID string) ([]domain.Deck, error) {
	decks := make([]domain.Deck, 0)
	rows, err := d.db.Query(`SELECT id, name, description, labels, vector_image_id, hidden, promo 
		FROM decks 
		WHERE !hidden OR (
		    hidden AND id IN (select deck_id FROM unlocked_decks WHERE client_id = ?)
		)
	`, clientID)
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

func (d *DecksMsStorage) SaveDeck(deck domain.Deck) error {
	return d.SaveDecks([]domain.Deck{deck})
}

func (d *DecksMsStorage) SaveDecks(decks []domain.Deck) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	for _, deck := range decks {
		_, err := tx.Exec(fmt.Sprintf(
			"INSERT INTO decks (%v) VALUES (%v)",
			"id, name, description, labels, vector_image_id, hidden, promo, language_code",
			"?,?,?,?,?,?,?,?",
		),
			deck.ID, deck.Name, deck.Description, strings.Join(deck.Labels, ";"), deck.Image, deck.IsHidden,
			deck.PromoCode, "RU",
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

func (d *DecksMsStorage) UnlockDeck(clientID, deckID string) error {
	_, err := d.db.Exec("INSERT INTO unlocked_decks (client_id, deck_id) VALUES (?, ?)",
		clientID, deckID)
	return err
}

func NewDecksMsStorage(
	db *sql.DB,
) *DecksMsStorage {
	return &DecksMsStorage{
		db: db,
	}
}
