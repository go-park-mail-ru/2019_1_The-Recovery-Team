package rest

import (
	"testing"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/infrastructure/repository/memory/game"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/infrastructure/repository/postgresql/profile"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/infrastructure/repository/redis/session"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/usecase"
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
