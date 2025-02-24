package adapters

import (
	"dc_honest/src/internal/adapters"
	"dc_honest/src/internal/core"
	"dc_honest/src/internal/infrastructure/mock"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func initShuffleAdapter(t *testing.T) (
	*core.Config,
	*adapters.ShuffleHttpAdapter,
	*gin.Engine,
) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal("Error loading .env file")
	}
	config := core.NewConfig()
	service := mock.NewShuffleMock()
	engine := gin.Default()
	adapter := adapters.NewShuffleHttpAdapter(engine, service)

	return config, adapter, engine
}

func TestShuffleHttpAdapter(t *testing.T) {
	_, _, engine := initShuffleAdapter(t)

	type TestCase struct {
		Name               string
		EntityID           string
		ClientID           string
		ExpectedStatusCode int
	}

	t.Run("Shuffle deck test", func(t *testing.T) {
		cases := []TestCase{
			{
				Name:               "Error if client not uuid",
				EntityID:           "12345",
				ClientID:           "some",
				ExpectedStatusCode: http.StatusBadRequest,
			},
			{
				Name:               "Ok if client is valid uuid",
				EntityID:           "12345",
				ClientID:           "834f506e-7a28-400e-8945-ffef09eafd65",
				ExpectedStatusCode: http.StatusOK,
			},
		}

		for _, tc := range cases {
			t.Run(tc.Name, func(t *testing.T) {
				r := httptest.NewRequest("POST", "/v1/decks/"+tc.EntityID+"/shuffle?clientId="+tc.ClientID, nil)
				w := httptest.NewRecorder()
				engine.ServeHTTP(w, r)

				assert.Equal(t, tc.ExpectedStatusCode, w.Code)
			})
		}
	})

	t.Run("Shuffle level test", func(t *testing.T) {
		cases := []TestCase{
			{
				Name:               "Error if client not uuid",
				EntityID:           "12345",
				ClientID:           "some",
				ExpectedStatusCode: http.StatusBadRequest,
			},
			{
				Name:               "Ok if client is valid uuid",
				EntityID:           "12345",
				ClientID:           "834f506e-7a28-400e-8945-ffef09eafd65",
				ExpectedStatusCode: http.StatusOK,
			},
		}

		for _, tc := range cases {
			t.Run(tc.Name, func(t *testing.T) {
				r := httptest.NewRequest("POST", "/v1/levels/"+tc.EntityID+"/shuffle?clientId="+tc.ClientID, nil)
				w := httptest.NewRecorder()
				engine.ServeHTTP(w, r)

				assert.Equal(t, tc.ExpectedStatusCode, w.Code)
			})
		}
	})
}
