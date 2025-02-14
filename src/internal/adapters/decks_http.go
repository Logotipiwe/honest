package adapters

import (
	"dc_honest/src/internal/core/ports/input"
	"dc_honest/src/pkg"
	"github.com/gin-gonic/gin"
)

type DecksAdapterHttp struct {
	service input.DecksPort
}

func NewDecksAdapterHttp(
	router *gin.Engine,
	service input.DecksPort,
) *DecksAdapterHttp {
	d := &DecksAdapterHttp{
		service: service,
	}

	router.GET("/ping", pkg.WithError(d.Ping))
	router.GET("/v1/decks", pkg.WithError(d.Decks))

	return d
}

// Ping godoc
// @Summary      Do ping
// @Description  Do ping desc
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Success      200  {object}  string
// @Failure      500  {object}  string
// @Router       /ping [get]
func (d *DecksAdapterHttp) Ping(c *gin.Context) error {
	c.JSON(200, gin.H{"ok": true})
	return nil
}

// Decks godoc
// @Summary      Get all public decks
// @Accept       json
// @Produce      json
// @Param		 client_id query string true "client id"
// @Success      200  {object}  DecksAnswer
// @Router       /v1/decks [get]
func (d *DecksAdapterHttp) Decks(c *gin.Context) error {
	clientID := c.Query("client_id")
	decks, err := d.service.GetDecksForMainPage(clientID)
	if err != nil {
		return err
	}
	dtos := make([]DeckDto, len(decks))
	for i, deck := range decks {
		dtos[i] = ToDto(deck)
	}
	c.JSON(200, DecksAnswer{
		Ok:    true,
		Decks: dtos,
	})
	return nil
}
