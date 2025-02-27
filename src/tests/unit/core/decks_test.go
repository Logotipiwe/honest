package core

import (
	"dc_honest/src/internal/core/domain"
	"dc_honest/src/internal/core/service"
	"dc_honest/src/internal/infrastructure/mock"
	. "dc_honest/src/pkg"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDecksForMainPage(t *testing.T) {
	storage := mock.NewDecksStorageMock()
	decksService := service.NewDecksService(storage)

	type TestCase struct {
		Name     string
		ClientID string
		Decks    []domain.Deck
		Expected []domain.Deck
	}

	cases := []TestCase{
		{
			"Works",
			"1",
			[]domain.Deck{
				{
					ID:          "1",
					Name:        "name1",
					Description: "",
					Labels:      nil,
					Image:       "",
					IsHidden:    false,
				},
				{
					ID:          "2",
					Name:        "",
					Description: "",
					Labels:      []string{"l1", "l2"},
					Image:       "im",
					IsHidden:    false,
					PromoCode:   P("promo"),
				},
			},
			[]domain.Deck{
				{
					ID:          "1",
					Name:        "name1",
					Description: "",
					Labels:      nil,
					Image:       "",
					IsHidden:    false,
				},
				{
					ID:          "2",
					Name:        "",
					Description: "",
					Labels:      []string{"l1", "l2"},
					Image:       "im",
					IsHidden:    false,
					PromoCode:   P("promo"),
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			defer storage.Clean()
			storage.SetDecks(tc.Decks)
			result, err := decksService.GetAvailableDecks(tc.ClientID)
			assert.Nil(t, err)
			assert.Equal(t, len(tc.Expected), len(result))
			for i, d := range tc.Expected {
				assert.EqualExportedValues(t, d, result[i])
			}
		})
	}
}
