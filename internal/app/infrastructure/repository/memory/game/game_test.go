package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewGameRepo(t *testing.T) {
	log, _ := zap.NewProduction()
	assert.NotEmpty(t, NewGameRepo(log),
		"Doesn't create game repository instance")
}
