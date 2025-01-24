package adapters

import "dc_honest/src/internal/core/domain"

type DecksAnswer struct {
	Ok    bool      `json:"ok"`
	Decks []DeckDto `json:"decks"`
}

type DeckDto struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Labels      []string `json:"labels"`
	ImageID     string   `json:"image_id"`
}

func ToDto(d domain.Deck) DeckDto {
	return DeckDto{
		ID:          d.ID,
		Name:        d.Name,
		Description: d.Description,
		Labels:      d.Labels,
		ImageID:     d.Image,
	}
}
