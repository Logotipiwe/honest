package tests

import (
	"dc_honest/src/internal/adapters/output"
	"dc_honest/src/internal/core/domain"
	"dc_honest/src/internal/core/service"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecks(t *testing.T) {
	t.Run("GetDecksForMainPage gets", func(t *testing.T) {
		storageMock := output.NewDecksStorageMock()
		d := service.NewDecksService(storageMock)
		saved := []domain.Deck{
			{
				ID:          "1",
				Name:        "name1",
				Description: nil,
				Labels:      nil,
				Image:       "",
				IsHidden:    false,
				PromoCode:   "",
			},
			{
				ID:          "2",
				Name:        "",
				Description: nil,
				Labels:      []string{"l1", "l2"},
				Image:       "im",
				IsHidden:    false,
				PromoCode:   "promo",
			},
		}
		storageMock.SetDecksForClient("1", saved)

		decks, err := d.GetDecksForMainPage("1")
		assert.Nil(t, err)
		assert.Equal(t, len(saved), len(decks))
	})
}
