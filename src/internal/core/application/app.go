package application

import (
	"dc_honest/src/internal/core/ports/input"
)

type App struct {
	input.DecksPort
	input.ShufflePort
}
