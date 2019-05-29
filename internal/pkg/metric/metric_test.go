package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const namespace = "test"

func TestRegisterAccessHitsMetric(t *testing.T) {
	RegisterAccessHitsMetric(namespace)
	assert.NotEmpty(t, AccessHits, "Doesn't register metric")
}

func TestRegisterTotalRoomsMetric(t *testing.T) {
	RegisterTotalRoomsMetric(namespace)
	assert.NotEmpty(t, TotalRooms, "Doesn't register metric")
}

func TestRegisterTotalPlayersMetric(t *testing.T) {
	RegisterTotalPlayersMetric(namespace)
	assert.NotEmpty(t, TotalPlayers, "Doesn't register metric")
}
