package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/infrastructure/repository/memory/game"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/infrastructure/repository/postgresql/profile"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const (
	gameUrl = "ws://127.0.0.1:8080/api/v1/game.ws"
)

var testCaseSearch = []struct {
	name       string
	profileId  uint64
	sessionId  string
	statusCode int
}{
	{
		name:       "Test with unauthorized user",
		profileId:  profile.ForbiddenProfileId,
		sessionId:  "",
		statusCode: http.StatusUnauthorized,
	},
	{
		name:       "Test with correct data",
		profileId:  profile.ExistingProfileId,
		sessionId:  "",
		statusCode: http.StatusBadRequest,
	},
}

func TestSearch(t *testing.T) {
	log, _ := zap.NewProduction()

	for _, testCase := range testCaseSearch {
		t.Run(testCase.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", gameUrl, nil)
			ctx := context.WithValue(req.Context(), "logger", log)
			ctx = context.WithValue(ctx, ProfileID, testCase.profileId)
			ctx = context.WithValue(ctx, SessionID, testCase.sessionId)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()

			profileInteractor := usecase.NewProfileInteractor(&profile.ProfileRepoMock{})
			gameInteractor := usecase.NewGameInteractor(&game.GameRepoMock{})
			Search(profileInteractor, gameInteractor)(w, req, nil)

			assert.Equal(t, testCase.statusCode, w.Code,
				"Wrong status code")
		})
	}
}
