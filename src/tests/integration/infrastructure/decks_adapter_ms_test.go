package infrastructure

import (
	"database/sql"
	"dc_honest/src/internal/core"
	"dc_honest/src/internal/core/domain"
	"dc_honest/src/internal/infrastructure/ms"
	"dc_honest/src/internal/infrastructure/ms/flyway"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"testing"
)

func p(v bool) *bool {
	return &v
}

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
		Expected      []domain.Deck
		ExpectSaveErr *bool
	}

	testCases := []TestCase{
		{
			Name: "Promo can't duplicate",
			Decks: []domain.Deck{
				{ID: "1", Name: "name 1", PromoCode: "c"},
				{ID: "2", Name: "name 2", PromoCode: "c"},
			},
			Expected:      []domain.Deck{},
			ExpectSaveErr: p(true),
		},
		{
			Name: "Doesn't show hidden decks",
			Decks: []domain.Deck{
				{ID: "1", Name: "name 1", PromoCode: "c"},
				{ID: "2", Name: "name 2", IsHidden: true},
			},
			Expected: []domain.Deck{
				{ID: "1", Name: "name 1", PromoCode: "c"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			defer fw.Clean()
			assert.Nil(t, fw.Migrate())
			err := storage.SaveDecks(tc.Decks)
			if tc.ExpectSaveErr != nil && *tc.ExpectSaveErr {
				assert.Error(t, err)
				return
			} else {
				if !assert.Nil(t, err) {
					return
				}
			}
			decks, err := storage.GetDecksForClient("1")
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
