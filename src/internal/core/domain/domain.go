package domain

type Deck struct {
	ID          string
	Name        string
	Description string
	Labels      []string
	Image       string
	IsHidden    bool
	PromoCode   string
}

type Level struct {
	ID          string
	Deck        Deck
	Name        string
	Order       int
	Description *string
	Color       string
}

type Question struct {
	ID             string
	Level          Level
	Text           string
	AdditionalText *string
}
