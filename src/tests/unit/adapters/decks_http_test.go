package adapters

import (
	"dc_honest/src/internal/adapters"
	"dc_honest/src/internal/core/domain"
	. "dc_honest/src/internal/core/service"
	"dc_honest/src/internal/infrastructure/mock"
	. "dc_honest/src/pkg"
	"dc_honest/src/tests"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func readDtos(t *testing.T, w *httptest.ResponseRecorder) []adapters.DeckOutput {
	var dtos []adapters.DeckOutput
	err := json.Unmarshal(w.Body.Bytes(), &dtos)
	assert.Nil(t, err)
	return dtos
}

func TestDecksHttpAdapter(t *testing.T) {
	tests.LoadTestEnv(t)
	storage := mock.NewDecksStorageMock()
	service := NewDecksService(storage)
	engine := gin.Default()

	adapters.NewDecksAdapterHttp(engine, service)

	type TestCase struct {
		Name             string
		ClientID         string
		DecksReturned    []domain.Deck
		ExpectedCalls    [][]any
		UseCamelClientID *bool
		Expected         *[]adapters.DeckOutput
	}

	cases := []TestCase{
		{Name: "Works", ClientID: "1",
			DecksReturned: []domain.Deck{
				{
					ID:          "1",
					Name:        "n1",
					Description: "desc",
					Labels:      nil,
					Image:       "image1",
					IsHidden:    false,
				},
				{
					ID:          "2",
					Name:        "n2",
					Description: "",
					Labels:      []string{"l1", "l2"},
					Image:       "image2",
					IsHidden:    false,
				},
			},
			ExpectedCalls: [][]any{
				{"1"},
			},
			Expected: &[]adapters.DeckOutput{
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
		{
			Name:             "Works with camel clientId param",
			ClientID:         "4",
			UseCamelClientID: P(true),
			ExpectedCalls: [][]any{
				{"4"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			defer storage.Clean()
			storage.SetDecks(tc.DecksReturned)
			var clientIdKey string
			if tc.UseCamelClientID != nil && *tc.UseCamelClientID {
				clientIdKey = "clientId"
			} else {
				clientIdKey = "client_id"
			}
			r := httptest.NewRequest("GET", "/v1/decks?"+clientIdKey+"="+tc.ClientID, nil)
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, r)

			assert.Equal(t, tc.ExpectedCalls, storage.GetAvailableDecksCalls)

			assert.Equal(t, 200, w.Code)
			dtos := readDtos(t, w)
			if tc.Expected != nil {
				assert.Equal(t, len(*tc.Expected), len(dtos))
				assert.True(t, SameElements(dtos, *tc.Expected))
			}
		})
	}
}
