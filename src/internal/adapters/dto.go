package adapters

import "dc_honest/src/internal/core/domain"

type DeckOutput struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Labels      []string `json:"labels"`
	ImageID     string   `json:"image_id"`
}

func ToOutput(d domain.Deck) DeckOutput {
	return DeckOutput{
		ID:          d.ID,
		Name:        d.Name,
		Description: d.Description,
		Labels:      d.Labels,
		ImageID:     d.Image,
	}
}

type QuestionOutput struct {
	ID             string  `json:"id"`
	Text           string  `json:"text"`
	LevelID        string  `json:"level_id"`
	AdditionalText *string `json:"additional_text"`
	IsLast         bool    `json:"is_last"`
}

func ToOutputDto(q domain.Question, isLast bool) QuestionOutput {
	return QuestionOutput{
		ID:             q.ID,
		Text:           q.Text,
		LevelID:        q.Level.ID,
		AdditionalText: q.AdditionalText,
		IsLast:         isLast,
	}
}
