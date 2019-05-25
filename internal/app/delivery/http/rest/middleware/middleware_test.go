package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/pkg/metric"

	"go.uber.org/zap"

	sessionGrpc "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/session"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/infrastructure/repository/redis/session"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

const (
	authenticationUrl = "http://127.0.0.1:8080/api/v1/authentication/"
	corsUrl           = "http://127.0.0.1:8080/api/v1/cors/"
	recoveryUrl       = "http://127.0.0.1:8080/api/v1/recovery/"
)

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

			sessionManager := sessionGrpc.NewClientMock()
			Authentication(&sessionManager, fakeHandler)(w, req, nil)

			assert.Equal(t, testCase.statusCode, w.Code,
				"Wrong status code")
		})
	}
}

var testCaseSessionMiddleware = []struct {
	name       string
	token      string
	hasCookie  bool
	statusCode int
}{
	{
		name:       "Test without cookie",
		hasCookie:  false,
		statusCode: http.StatusOK,
	},
	{
		name:       "Test with not existing cookie",
		token:      session.Unauthorized,
		hasCookie:  true,
		statusCode: http.StatusOK,
	},
	{
		name:       "Test with correct cookie",
		token:      session.Authorized,
		hasCookie:  true,
		statusCode: http.StatusOK,
	},
}

func TestSessionMiddleware(t *testing.T) {
	for _, testCase := range testCaseSessionMiddleware {
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

			sessionManager := sessionGrpc.NewClientMock()
			SessionMiddleware(&sessionManager, fakeHandler)(w, req, nil)

			assert.Equal(t, testCase.statusCode, w.Code,
				"Wrong status code")
		})
	}
}

func TestCORSMiddleware(t *testing.T) {
	req := httptest.NewRequest("GET", corsUrl, nil)
	req.Header.Set("Origin", "https://sadislands.ru")
	w := httptest.NewRecorder()

	CORSMiddleware(fakeHandler)(w, req, nil)

	assert.Equal(t, http.StatusOK, w.Code,
		"Wrong status code")
}

func TestRecoveryMiddleware(t *testing.T) {
	req := httptest.NewRequest("GET", recoveryUrl, nil)
	w := httptest.NewRecorder()

	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	RecoverMiddleware(log, fakeHandler)(w, req, nil)

	assert.Equal(t, http.StatusOK, w.Code,
		"Wrong status code")
}

func TestLoggerMiddleware(t *testing.T) {
	req := httptest.NewRequest("GET", recoveryUrl, nil)
	w := httptest.NewRecorder()

	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	LoggerMiddleware(log, fakeHandler)(w, req, nil)

	assert.Equal(t, http.StatusOK, w.Code,
		"Wrong status code")
}

func TestAccessHitsMiddleware(t *testing.T) {
	req := httptest.NewRequest("GET", recoveryUrl, nil)
	w := httptest.NewRecorder()
	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	ctx := context.WithValue(context.Background(), "logger", log)

	metric.RegisterAccessHitsMetric("test")
	AccessHitsMiddleware(fakeHandler)(w, req.WithContext(ctx), nil)

	assert.Equal(t, http.StatusOK, w.Code,
		"Wrong status code")
}
