package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	profileGrpc "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/profile"

	"go.uber.org/zap"

	sessionGrpc "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/session"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/http/rest/middleware"
	profileDomain "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/domain/profile"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/infrastructure/repository/postgresql/profile"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/infrastructure/repository/redis/session"

	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
)

const (
	ProfileID = 0
	SessionID = 1
)

const (
	sessionUrl = "http://127.0.0.1:8080/api/v1/session/"
)

var testCaseGetSession = []struct {
	name       string
	id         uint64
	statusCode int
	body       profileDomain.ID
}{
	{
		name:       "Test correct getting session",
		id:         1,
		statusCode: http.StatusOK,
		body:       profileDomain.ID{Id: 1},
	},
}

var testCasePostSession = []struct {
	name       string
	email      string
	password   string
	setCookie  bool
	statusCode int
}{
	{
		name:       "Test without request data",
		setCookie:  false,
		statusCode: http.StatusBadRequest,
	},
	{
		name:       "Test with invalid data",
		email:      profile.IncorrectProfileEmail,
		setCookie:  false,
		statusCode: http.StatusUnprocessableEntity,
	},
	{
		name:       "Test with invalid password",
		email:      profile.ExistingProfileEmail,
		password:   profile.ForbiddenProfilePassword,
		setCookie:  false,
		statusCode: http.StatusUnprocessableEntity,
	},
	{
		name:       "Test with not existing email",
		email:      profile.NotExistingProfileEmail,
		password:   profile.ExistingProfilePassword,
		setCookie:  false,
		statusCode: http.StatusInternalServerError,
	},
	{
		name:       "Test with setting session error",
		email:      profile.ForbiddenProfileEmail,
		password:   profile.ExistingProfilePassword,
		setCookie:  false,
		statusCode: http.StatusInternalServerError,
	},
	{
		name:       "Test with correct data",
		email:      profile.ExistingProfileEmail,
		password:   profile.ExistingProfilePassword,
		setCookie:  true,
		statusCode: http.StatusOK,
	},
}

var testCaseDeleteSession = []struct {
	name         string
	token        string
	deleteCookie bool
	statusCode   int
}{
	{
		name:         "Test delete authorized",
		token:        session.Authorized,
		deleteCookie: true,
		statusCode:   http.StatusOK,
	},
	{
		name:         "Test delete unauthorized",
		token:        session.Unauthorized,
		deleteCookie: false,
		statusCode:   http.StatusNotFound,
	},
}

func TestGetSession(t *testing.T) {
	for _, testCase := range testCaseGetSession {
		t.Run(testCase.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", sessionUrl, nil)
			ctx := context.WithValue(req.Context(), middleware.ProfileID, testCase.id)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()

			GetSession()(w, req, nil)

			assert.Equal(t, testCase.statusCode, w.Code,
				"Wrong status code")

			expectedBody, _ := easyjson.Marshal(testCase.body)
			assert.Equal(t, string(expectedBody), w.Body.String(),
				"Wrong response body")
		})
	}
}

func TestPostSession(t *testing.T) {
	log, _ := zap.NewProduction()

	for _, testCase := range testCasePostSession {
		t.Run(testCase.name, func(t *testing.T) {
			var req *http.Request
			if testCase.email != "" || testCase.password != "" {
				bodyRaw, _ := easyjson.Marshal(&profileDomain.Login{
					Email:    testCase.email,
					Password: testCase.password,
				})
				req = httptest.NewRequest("POST", sessionUrl, bytes.NewReader(bodyRaw))
			} else {
				req = httptest.NewRequest("POST", sessionUrl, nil)
			}

			ctx := context.WithValue(req.Context(), "logger", log)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()

			sessionManager := sessionGrpc.NewClientMock()
			profileManager := profileGrpc.NewClientMock()
			PostSession(&profileManager, &sessionManager)(w, req, nil)

			assert.Equal(t, testCase.statusCode, w.Code,
				"Wrong status code")

			cookies := w.Result().Cookies()
			if !testCase.setCookie {
				assert.Empty(t, cookies, "Sets cookie on wrong data")
				return
			}

			assert.NotEmpty(t, cookies, "Doesn't set cookie on correct data")

			expectedBody, _ := easyjson.Marshal(profileDomain.Profile{
				Info: profileDomain.Info{
					ID: profile.DefaultProfileId,
				},
			})

			assert.Equal(t, string(expectedBody), w.Body.String(),
				"Wrong response body")
		})
	}
}

func TestDeleteSession(t *testing.T) {
	log, _ := zap.NewProduction()

	for _, testCase := range testCaseDeleteSession {
		t.Run(testCase.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", sessionUrl, nil)
			ctx := context.WithValue(req.Context(), middleware.SessionID, testCase.token)
			ctx = context.WithValue(ctx, "logger", log)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()

			sessionManager := sessionGrpc.NewClientMock()
			DeleteSession(&sessionManager)(w, req, nil)

			assert.Equal(t, testCase.statusCode, w.Code,
				"Wrong status code")

			cookies := w.Result().Cookies()
			if testCase.deleteCookie {
				assert.NotEqual(t, nil, cookies, "Doesn't delete cookie")
				return
			}

			assert.Empty(t, cookies, "Deletes cookie on wrong token")
		})
	}
}
