package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/infrastructure/repository/redis/session"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/usecase"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

const authenticationUrl = "http://127.0.0.1:8080/api/v1/authentication/"

func fakeHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
}

var testCaseAuthentication = []struct {
	name       string
	token      string
	hasCookie  bool
	statusCode int
}{
	{
		name:       "Test without cookie",
		hasCookie:  false,
		statusCode: http.StatusUnauthorized,
	},
	{
		name:       "Test with not existing cookie",
		token:      session.Unauthorized,
		hasCookie:  true,
		statusCode: http.StatusUnauthorized,
	},
	{
		name:       "Test with correct cookie",
		token:      session.Authorized,
		hasCookie:  true,
		statusCode: http.StatusOK,
	},
}

func TestAuthentication(t *testing.T) {
	for _, testCase := range testCaseAuthentication {
		t.Run(testCase.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", authenticationUrl, nil)
			if testCase.hasCookie {
				req.AddCookie(&http.Cookie{
					Name:     "session_id",
					Value:    testCase.token,
					Path:     "/",
					HttpOnly: true,
				})
			}

			w := httptest.NewRecorder()

			sessionInteractor := usecase.NewSessionInteractor(&session.SessionRepoMock{})
			Authentication(sessionInteractor, fakeHandler)(w, req, nil)

			assert.Equal(t, testCase.statusCode, w.Code,
				"Wrong status code")
		})
	}
}
