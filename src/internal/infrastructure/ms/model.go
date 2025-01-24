package ms

import (
	"dc_honest/src/internal/core/domain"
	"strings"
)

type DeckModel struct {
	ID          string
	Name        string
	Description string
	Labels      string
	Image       string
	IsHidden    bool
	PromoCode   string
}

func (m DeckModel) ToDeck() domain.Deck {
	var labels []string = nil
	if m.Labels != "" {
		labels = strings.Split(m.Labels, ",")
	}
	return domain.Deck{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Labels:      labels,
		Image:       m.Image,
		IsHidden:    m.IsHidden,
		PromoCode:   m.PromoCode,
	}
}

type LevelModel struct {
	ID          string
	DeckID      string
	Name        string
	Description *string
	Color       string
}

type QuestionModel struct {
	ID      string
	LevelID string
	Text    string
}
