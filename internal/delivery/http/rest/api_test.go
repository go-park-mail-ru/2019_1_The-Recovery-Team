package rest

import (
	"sadislands/internal/infrastructure/repository/memory/game"
	"sadislands/internal/infrastructure/repository/postgresql/profile"
	"sadislands/internal/infrastructure/repository/redis/session"
	"sadislands/internal/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewRestApi(t *testing.T) {
	log, _ := zap.NewProduction()
	profileInteractor := usecase.NewProfileInteractor(&profile.ProfileRepoMock{})
	sessionInteractor := usecase.NewSessionInteractor(&session.SessionRepoMock{})
	gameInteractor := usecase.NewGameInteractor(&game.GameRepoMock{})

	api := NewRestApi(profileInteractor, sessionInteractor, gameInteractor, log)
	assert.NotEmpty(t, api,
		"Returns empty api router")
}
