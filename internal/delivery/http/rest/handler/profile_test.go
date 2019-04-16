package handler

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	profileDomain "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/domain/profile"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/infrastructure/repository/postgresql/profile"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/infrastructure/repository/redis/session"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/usecase"

	"github.com/julienschmidt/httprouter"

	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const (
	profileUrl = "http://127.0.0.1:8080/api/v1/profiles"
	scoreUrl   = "http://127.0.0.1:8080/api/v1/scores"
)

var testCaseGetProfile = []struct {
	name       string
	id         string
	profileId  string
	statusCode int
}{
	{
		name:       "Test with incorrect id",
		id:         "-1",
		profileId:  profile.DefaultProfileIdStr,
		statusCode: http.StatusNotFound,
	},
	{
		name:       "Test with not existing profile",
		id:         profile.NotExistingProfileIdStr,
		profileId:  profile.DefaultProfileIdStr,
		statusCode: http.StatusNotFound,
	},
	{
		name:       "Test with error in database",
		id:         profile.ForbiddenProfileIdStr,
		profileId:  profile.DefaultProfileIdStr,
		statusCode: http.StatusInternalServerError,
	},
	{
		name:       "Test getting foreign profile",
		id:         profile.ExistingProfileIdStr,
		profileId:  profile.DefaultProfileIdStr,
		statusCode: http.StatusOK,
	},
	{
		name:       "Test getting own profile",
		id:         profile.ExistingProfileIdStr,
		profileId:  profile.ExistingProfileIdStr,
		statusCode: http.StatusOK,
	},
}

var testCaseGetProfiles = []struct {
	name       string
	email      string
	nickname   string
	statusCode int
}{
	{
		name:       "Test not existing email",
		email:      profile.NotExistingProfileEmail,
		statusCode: http.StatusNotFound,
	},
	{
		name:       "Test database error on forbidden email",
		email:      profile.ForbiddenProfileEmail,
		statusCode: http.StatusInternalServerError,
	},
	{
		name:       "Test existing email",
		email:      profile.ExistingProfileEmail,
		statusCode: http.StatusOK,
	},
	{
		name:       "Test not existing nickname",
		nickname:   profile.NotExistingProfileNickname,
		statusCode: http.StatusNotFound,
	},
	{
		name:       "Test database error on forbidden nickname",
		nickname:   profile.ForbiddenProfileNickname,
		statusCode: http.StatusInternalServerError,
	},
	{
		name:       "Test existing nickname",
		nickname:   profile.ExistingProfileNickname,
		statusCode: http.StatusOK,
	},
	{
		name:       "Test with invalid email and password",
		statusCode: http.StatusBadRequest,
	},
}

var testCasePutProfile = []struct {
	name       string
	id         string
	profileId  string
	email      string
	nickname   string
	statusCode int
}{
	{
		name:       "Test with invalid profile id",
		id:         "-1",
		profileId:  profile.NotExistingProfileIdStr,
		email:      profile.ExistingProfileEmail,
		nickname:   profile.ExistingProfileNickname,
		statusCode: http.StatusNotFound,
	},
	{
		name:       "Test without permissions",
		id:         profile.NotExistingProfileIdStr,
		profileId:  profile.ExistingProfileIdStr,
		email:      profile.ExistingProfileEmail,
		nickname:   profile.ExistingProfileNickname,
		statusCode: http.StatusForbidden,
	},
	{
		name:       "Test without data",
		id:         profile.ExistingProfileIdStr,
		profileId:  profile.ExistingProfileIdStr,
		statusCode: http.StatusBadRequest,
	},
	{
		name:       "Test with invalid update data",
		id:         profile.ExistingProfileIdStr,
		profileId:  profile.ExistingProfileIdStr,
		email:      profile.IncorrectProfileEmail,
		nickname:   profile.ExistingProfileNickname,
		statusCode: http.StatusUnprocessableEntity,
	},
	{
		name:       "Test with conflict data",
		id:         profile.ExistingProfileIdStr,
		profileId:  profile.ExistingProfileIdStr,
		email:      profile.ExistingProfileEmail,
		nickname:   profile.ExistingProfileNickname,
		statusCode: http.StatusConflict,
	},
	{
		name:       "Test with database error",
		id:         profile.ExistingProfileIdStr,
		profileId:  profile.ExistingProfileIdStr,
		email:      profile.ForbiddenProfileEmail,
		nickname:   profile.ExistingProfileNickname,
		statusCode: http.StatusInternalServerError,
	},
	{
		name:       "Test with correct data",
		id:         profile.ExistingProfileIdStr,
		profileId:  profile.ExistingProfileIdStr,
		email:      profile.NotExistingProfileEmail,
		nickname:   profile.NotExistingProfileNickname,
		statusCode: http.StatusNoContent,
	},
}

var testPutProfilePassword = []struct {
	name        string
	id          string
	profileId   string
	password    string
	passwordOld string
	statusCode  int
}{
	{
		name:        "Test with invalid profile id",
		id:          "-1",
		profileId:   profile.ExistingProfileIdStr,
		password:    profile.NotExistingProfilePassword,
		passwordOld: profile.ExistingProfilePassword,
		statusCode:  http.StatusNotFound,
	},
	{
		name:        "Test without permissions",
		id:          profile.NotExistingProfileIdStr,
		profileId:   profile.ExistingProfileIdStr,
		password:    profile.NotExistingProfilePassword,
		passwordOld: profile.ExistingProfilePassword,
		statusCode:  http.StatusForbidden,
	},
	{
		name:       "Test without data",
		id:         profile.ExistingProfileIdStr,
		profileId:  profile.ExistingProfileIdStr,
		statusCode: http.StatusBadRequest,
	},
	{
		name:        "Test with invalid password",
		id:          profile.ExistingProfileIdStr,
		profileId:   profile.ExistingProfileIdStr,
		password:    profile.InvalidProfilePassword,
		passwordOld: profile.InvalidProfilePassword,
		statusCode:  http.StatusUnprocessableEntity,
	},
	{
		name:        "Test with incorrect old password",
		id:          profile.ExistingProfileIdStr,
		profileId:   profile.ExistingProfileIdStr,
		password:    profile.NotExistingProfilePassword,
		passwordOld: profile.NotExistingProfilePassword,
		statusCode:  http.StatusUnprocessableEntity,
	},
	{
		name:        "Test with database error",
		id:          profile.ExistingProfileIdStr,
		profileId:   profile.ExistingProfileIdStr,
		password:    profile.ForbiddenProfilePassword,
		passwordOld: profile.ForbiddenProfilePassword,
		statusCode:  http.StatusInternalServerError,
	},
	{
		name:        "Test with correct data",
		id:          profile.ExistingProfileIdStr,
		profileId:   profile.ExistingProfileIdStr,
		password:    profile.NotExistingProfilePassword,
		passwordOld: profile.ExistingProfilePassword,
		statusCode:  http.StatusNoContent,
	},
}

var testCasePostProfile = []struct {
	name       string
	email      string
	nickname   string
	password   string
	multipart  bool
	statusCode int
}{
	{
		name:       "Test not multipart",
		email:      profile.ExistingProfileEmail,
		nickname:   profile.ExistingProfileNickname,
		password:   profile.ExistingProfilePassword,
		multipart:  false,
		statusCode: http.StatusBadRequest,
	},
	{
		name:       "Test with incorrect data",
		multipart:  true,
		statusCode: http.StatusBadRequest,
	},
	{
		name:       "Test with invalid data",
		email:      profile.IncorrectProfileEmail,
		nickname:   profile.ExistingProfileNickname,
		password:   profile.InvalidProfilePassword,
		multipart:  true,
		statusCode: http.StatusUnprocessableEntity,
	},
	{
		name:       "Test with conflict data",
		email:      profile.ExistingProfileEmail,
		nickname:   profile.ExistingProfileNickname,
		password:   profile.ExistingProfilePassword,
		multipart:  true,
		statusCode: http.StatusConflict,
	},
	{
		name:       "Test with database error",
		email:      profile.ForbiddenProfileEmail,
		nickname:   profile.ExistingProfileNickname,
		password:   profile.ExistingProfilePassword,
		multipart:  true,
		statusCode: http.StatusInternalServerError,
	},
	{
		name:       "Test with redis error",
		email:      profile.NotExistingProfileEmail,
		nickname:   profile.NotExistingProfileNickname,
		password:   profile.IncorrectProfilePassword,
		multipart:  true,
		statusCode: http.StatusInternalServerError,
	},
	{
		name:       "Test with correct data",
		email:      profile.NotExistingProfileEmail,
		nickname:   profile.NotExistingProfileNickname,
		password:   profile.NotExistingProfilePassword,
		multipart:  true,
		statusCode: http.StatusCreated,
	},
}

var testCaseGetScoreboard = []struct {
	name       string
	limit      string
	start      string
	statusCode int
}{
	{
		name:       "Test with incorrect limit and start",
		limit:      "-",
		start:      "-",
		statusCode: http.StatusBadRequest,
	},
	{
		name:       "Test with database error",
		limit:      profile.ForbiddenLimitStr,
		start:      "-1",
		statusCode: http.StatusInternalServerError,
	},
	{
		name:       "Test with correct data",
		limit:      "100",
		start:      "-1",
		statusCode: http.StatusOK,
	},
}

func TestGetProfile(t *testing.T) {
	log, _ := zap.NewProduction()

	for _, testCase := range testCaseGetProfile {
		t.Run(testCase.name, func(t *testing.T) {
			profileId, _ := strconv.ParseUint(testCase.profileId, 10, 64)
			req := httptest.NewRequest("GET", profileUrl+"/"+testCase.id, nil)
			ctx := context.WithValue(req.Context(), "logger", log)
			ctx = context.WithValue(ctx, ProfileID, profileId)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()

			ps := httprouter.Params{
				{
					Key:   "id",
					Value: testCase.id,
				},
			}

			profileInteractor := usecase.NewProfileInteractor(&profile.ProfileRepoMock{})
			GetProfile(profileInteractor)(w, req, ps)

			assert.Equal(t, testCase.statusCode, w.Code,
				"Wrong status code")

			if testCase.statusCode == http.StatusOK {
				id, _ := strconv.ParseUint(testCase.id, 10, 64)

				expectedBody, _ := easyjson.Marshal(profileDomain.Profile{
					Info: profileDomain.Info{
						ID: id,
					},
				})

				assert.Equal(t, string(expectedBody), w.Body.String(),
					"Wrong response body")
			}
		})
	}
}

func TestGetProfiles(t *testing.T) {
	log, _ := zap.NewProduction()

	for _, testCase := range testCaseGetProfiles {
		t.Run(testCase.name, func(t *testing.T) {

			var url string

			if testCase.email != "" {
				url = profileUrl + "/?email=" + testCase.email
			} else if testCase.nickname != "" {
				url = profileUrl + "/?nickname=" + testCase.nickname
			} else {
				url = profileUrl
			}
			req := httptest.NewRequest("GET", url, nil)
			ctx := context.WithValue(req.Context(), "logger", log)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()

			profileInteractor := usecase.NewProfileInteractor(&profile.ProfileRepoMock{})
			GetProfiles(profileInteractor)(w, req, nil)

			assert.Equal(t, testCase.statusCode, w.Code,
				"Wrong status code")
		})
	}
}

func TestPutProfile(t *testing.T) {
	log, _ := zap.NewProduction()

	for _, testCase := range testCasePutProfile {
		t.Run(testCase.name, func(t *testing.T) {
			var req *http.Request
			if testCase.email != "" || testCase.nickname != "" {
				bodyRaw, _ := easyjson.Marshal(&profileDomain.UpdateInfo{
					Email:    testCase.email,
					Nickname: testCase.nickname,
				})
				req = httptest.NewRequest("POST", profileUrl+"/"+testCase.id, bytes.NewReader(bodyRaw))
			} else {
				req = httptest.NewRequest("POST", profileUrl+"/"+testCase.id, nil)
			}

			profileId, _ := strconv.ParseUint(testCase.profileId, 10, 64)
			ctx := context.WithValue(req.Context(), "logger", log)
			ctx = context.WithValue(ctx, ProfileID, profileId)

			w := httptest.NewRecorder()

			ps := httprouter.Params{
				{
					Key:   "id",
					Value: testCase.id,
				},
			}

			profileInteractor := usecase.NewProfileInteractor(&profile.ProfileRepoMock{})
			PutProfile(profileInteractor)(w, req.WithContext(ctx), ps)

			assert.Equal(t, testCase.statusCode, w.Code,
				"Wrong status code")
		})
	}
}

func TestPutProfilePassword(t *testing.T) {
	log, _ := zap.NewProduction()

	for _, testCase := range testPutProfilePassword {
		t.Run(testCase.name, func(t *testing.T) {
			var req *http.Request
			if testCase.password != "" || testCase.passwordOld != "" {
				bodyRaw, _ := easyjson.Marshal(&profileDomain.UpdatePassword{
					Password:    testCase.password,
					PasswordOld: testCase.passwordOld,
				})
				req = httptest.NewRequest("PUT", profileUrl+"/password"+testCase.id, bytes.NewReader(bodyRaw))
			} else {
				req = httptest.NewRequest("PUT", profileUrl+"/password"+testCase.id, nil)
			}

			profileId, _ := strconv.ParseUint(testCase.profileId, 10, 64)
			ctx := context.WithValue(req.Context(), "logger", log)
			ctx = context.WithValue(ctx, ProfileID, profileId)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()

			ps := httprouter.Params{
				{
					Key:   "id",
					Value: testCase.id,
				},
			}

			profileInteractor := usecase.NewProfileInteractor(&profile.ProfileRepoMock{})
			PutProfilePassword(profileInteractor)(w, req, ps)

			assert.Equal(t, testCase.statusCode, w.Code,
				"Wrong status code")
		})
	}
}

func TestPostProfile(t *testing.T) {
	log, _ := zap.NewProduction()

	for _, testCase := range testCasePostProfile {
		t.Run(testCase.name, func(t *testing.T) {
			var buf bytes.Buffer
			writter := multipart.NewWriter(&buf)

			if testCase.nickname != "" {
				writter.WriteField("nickname", testCase.nickname)
			}
			if testCase.email != "" {
				writter.WriteField("email", testCase.email)
			}
			if testCase.password != "" {
				writter.WriteField("password", testCase.password)
			}
			writter.Close()

			req := httptest.NewRequest("POST", profileUrl, bytes.NewReader(buf.Bytes()))
			ctx := context.WithValue(req.Context(), "logger", log)
			req = req.WithContext(ctx)
			if testCase.multipart {
				req.Header.Set("Content-Type", writter.FormDataContentType())
			}

			w := httptest.NewRecorder()
			profileInteractor := usecase.NewProfileInteractor(&profile.ProfileRepoMock{})
			sessionInteractor := usecase.NewSessionInteractor(&session.SessionRepoMock{})
			PostProfile(profileInteractor, sessionInteractor)(w, req, nil)

			assert.Equal(t, testCase.statusCode, w.Code,
				"Wrong status code")
		})
	}
}

func TestGetScoreboard(t *testing.T) {
	log, _ := zap.NewProduction()

	for _, testCase := range testCaseGetScoreboard {
		t.Run(testCase.name, func(t *testing.T) {
			url := scoreUrl
			if testCase.limit != "" {
				url += "?limit=" + testCase.limit
			}
			if testCase.start != "" {
				if testCase.limit != "" {
					url += "&start=" + testCase.start
				} else {
					url += "?start=" + testCase.start
				}
			}

			req := httptest.NewRequest("GET", url, nil)
			ctx := context.WithValue(req.Context(), "logger", log)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()

			profileInteractor := usecase.NewProfileInteractor(&profile.ProfileRepoMock{})
			GetScoreboard(profileInteractor)(w, req, nil)

			assert.Equal(t, testCase.statusCode, w.Code,
				"Wrong status code")
		})
	}
}
