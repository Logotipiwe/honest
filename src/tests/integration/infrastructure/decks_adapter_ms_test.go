package infrastructure

import (
	"database/sql"
	"dc_honest/src/internal/core"
	"dc_honest/src/internal/core/domain"
	"dc_honest/src/internal/infrastructure/ms"
	"dc_honest/src/internal/infrastructure/ms/flyway"
	. "dc_honest/src/pkg"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecksAdapter(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal("Error loading .env file")
	}
	config := core.NewConfig()
	db, err := sql.Open("mysql", config.GetMysqlConnectionStr())
	if err != nil {
		t.Fatal(err)
	}
	storage := ms.NewDecksMsStorage(db)
	fw := flyway.NewFlyway(db, "../../../../data/migrations")

	type TestCase struct {
		Name          string
		Decks         []domain.Deck
		UnlockedDecks *[]struct{ clientID, deckID string }
		ClientID      string
		Expected      []domain.Deck
		ExpectSaveErr *bool
	}

	testCases := []TestCase{
		{
			Name: "Promo can't duplicate",
			Decks: []domain.Deck{
				{ID: "1", Name: "name 1", PromoCode: P("c")},
				{ID: "2", Name: "name 2", PromoCode: P("c")},
			},
			ClientID:      "1",
			Expected:      []domain.Deck{},
			ExpectSaveErr: P(true),
		},
		{
			Name: "Doesn't show hidden decks",
			Decks: []domain.Deck{
				{ID: "1", Name: "name 1", PromoCode: P("c")},
				{ID: "2", Name: "name 2", IsHidden: true},
			},
			ClientID: "1",
			Expected: []domain.Deck{
				{ID: "1", Name: "name 1", PromoCode: P("c")},
			},
		},
		{
			Name: "Show unlocked hidden decks",
			Decks: []domain.Deck{
				{ID: "1", Name: "name 1", IsHidden: false},
				{ID: "2", Name: "name 2", IsHidden: true},
				{ID: "3", Name: "name 3", IsHidden: true},
				{ID: "4", Name: "name 4", IsHidden: true},
			},
			UnlockedDecks: &[]struct{ clientID, deckID string }{
				{clientID: "client_1", deckID: "3"},
				{clientID: "client_1", deckID: "4"},
				{clientID: "client_2", deckID: "2"},
			},
			ClientID: "client_1",
			Expected: []domain.Deck{
				{ID: "1", Name: "name 1", IsHidden: false},
				{ID: "3", Name: "name 3", IsHidden: true},
				{ID: "4", Name: "name 4", IsHidden: true},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			assert.Nil(t, fw.Migrate())
			defer fw.Clean()
			err := storage.SaveDecks(tc.Decks)
			if tc.ExpectSaveErr != nil && *tc.ExpectSaveErr {
				assert.Error(t, err)
				return
			} else {
				if !assert.Nil(t, err) {
					return
				}
			}
			if tc.UnlockedDecks != nil {
				for _, deckToUnlock := range *tc.UnlockedDecks {
					err := storage.UnlockDeck(deckToUnlock.clientID, deckToUnlock.deckID)
					if !assert.Nil(t, err) {
						return
					}
				}
			}
			decks, err := storage.GetAvailableDecks(tc.ClientID)
			if !assert.Nil(t, err) {
				return
			}
			assert.Equal(t, len(tc.Expected), len(decks))
			for i, deck := range tc.Expected {
				assert.EqualExportedValues(t, deck, decks[i])
			}
		})
	}
}
