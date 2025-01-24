package adapters

import (
	"dc_honest/src/internal/adapters"
	"dc_honest/src/internal/core"
	"dc_honest/src/internal/core/domain"
	. "dc_honest/src/internal/core/service"
	"dc_honest/src/internal/infrastructure/mock"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func readDtos(t *testing.T, w *httptest.ResponseRecorder) adapters.DecksAnswer {
	dtos := adapters.DecksAnswer{}
	err := json.Unmarshal(w.Body.Bytes(), &dtos)
	assert.Nil(t, err)
	return dtos
}

func initAdapter(t *testing.T) (
	*core.Config,
	*mock.DecksStorageMock,
	*adapters.DecksAdapterHttp,
	*gin.Engine,
) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatal("Error loading .env file")
	}
	config := core.NewConfig()
	storage := mock.NewDecksStorageMock()
	service := NewDecksService(storage)
	engine := gin.Default()
	adapter := adapters.NewDecksAdapterHttp(engine, service)

	return config, storage, adapter, engine
}

func cleanup(storage *mock.DecksStorageMock) {
	storage.Clean()
}

func TestDecksHttpAdapter(t *testing.T) {
	_, storage, _, engine := initAdapter(t)

	defer cleanup(storage)

	type TestCase struct {
		Name     string
		ClientID string
		Decks    []domain.Deck
		Expected []adapters.DeckDto
	}

	cases := []TestCase{
		{Name: "Works", ClientID: "1",
			Decks: []domain.Deck{
				{
					ID:          "1",
					Name:        "n1",
					Description: "desc",
					Labels:      nil,
					Image:       "image1",
					IsHidden:    false,
					PromoCode:   "",
				},
				{
					ID:          "2",
					Name:        "n2",
					Description: "",
					Labels:      []string{"l1", "l2"},
					Image:       "image2",
					IsHidden:    false,
					PromoCode:   "",
				},
			},
			Expected: []adapters.DeckDto{
				{
					ID:          "1",
					Name:        "n1",
					Description: "desc",
					Labels:      nil,
					ImageID:     "image1",
				},
				{
					ID:          "2",
					Name:        "n2",
					Description: "",
					Labels:      []string{"l1", "l2"},
					ImageID:     "image2",
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			storage.SetDecks(tc.Decks)
			r := httptest.NewRequest("GET", "/v1/decks?client_id="+tc.ClientID, nil)
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, r)

			assert.Equal(t, 200, w.Code)
			dtos := readDtos(t, w)
			assert.True(t, dtos.Ok)
			assert.Equal(t, len(tc.Decks), len(dtos.Decks))
		})
	}
}
