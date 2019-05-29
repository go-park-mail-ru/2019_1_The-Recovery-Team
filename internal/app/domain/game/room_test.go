package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewRoom(t *testing.T) {
	closed := make(chan *Room, 5)
	log, _ := zap.NewProduction()

	assert.NotEmpty(t, NewRoom(log, closed),
		"Doesn't create room repository instance")
	close(closed)
}
