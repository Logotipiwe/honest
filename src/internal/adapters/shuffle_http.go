package adapters

import (
	"dc_honest/src/internal/core/ports/input"
	"dc_honest/src/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type ShuffleHttpAdapter struct {
	service input.ShufflePort
}

func NewShuffleHttpAdapter(
	router *gin.Engine,
	service input.ShufflePort,
) *ShuffleHttpAdapter {
	a := &ShuffleHttpAdapter{
		service: service,
	}

	router.POST("/api/v1/decks/:deckId/shuffle", pkg.WithError(a.ShuffleDeck))
	router.POST("/api/v1/levels/:levelId/shuffle", pkg.WithError(a.shuffleLevel))

	return a
}

// ShuffleDeck godoc
// @Summary      Перемешать вопросы в колоде и начать сначала
// @Param 		 clientId query string true "Client id"
// @Param 		 deckId path string true "Deck id"
// @Produce      json
// @Success      200
// @Router       /api/v1/decks/{deckId}/shuffle [post]
func (a *ShuffleHttpAdapter) ShuffleDeck(ctx *gin.Context) error {
	deckID := ctx.Param("deckId")
	clientIdStr := ctx.Query("clientId")
	clientID, err := uuid.Parse(clientIdStr)
	if err != nil {
		return ctx.AbortWithError(http.StatusBadRequest, err)
	}
	return a.service.ShuffleDeck(deckID, clientID.String())
}

// ShuffleLevel godoc
// @Summary      Перемешать вопросы в уровне и начать сначала
// @Param 		 clientId query string true "Client id"
// @Param 		 levelId path string true "Level id"
// @Produce      json
// @Success      200
// @Router       /api/v1/levels/{levelId}/shuffle [post]
func (a *ShuffleHttpAdapter) shuffleLevel(ctx *gin.Context) error {
	levelID := ctx.Param("levelId")
	clientIdStr := ctx.Query("clientId")
	clientID, err := uuid.Parse(clientIdStr)
	if err != nil {
		return ctx.AbortWithError(http.StatusBadRequest, err)
	}
	return a.service.ShuffleLevel(levelID, clientID.String())
}
